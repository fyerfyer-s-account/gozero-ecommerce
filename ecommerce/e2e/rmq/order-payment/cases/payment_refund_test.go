package cases

import (
	"context"
	"database/sql"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/order-payment/helpers"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/order-payment/setup"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/stretchr/testify/require"
)

func (s *PaymentTestSuite)TestPaymentRefundFlow() {
    // Initialize test context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    testCtx, err := setup.NewTestContext()
    require.NoError(s.T(), err)
    defer testCtx.Close()

    // Load test event fixture
    var event struct {
        OrderNo      string    `json:"order_no"`
        PaymentNo    string    `json:"payment_no"`
        RefundNo     string    `json:"refund_no"`
        RefundAmount float64   `json:"refund_amount"`
        Reason       string    `json:"reason"`
        Timestamp    time.Time `json:"timestamp"`
    }
    require.NoError(s.T(), helpers.LoadFixture("payment_refund.json", &event))

    // Initialize test data - Order
    order := &model.Orders{
        OrderNo:     event.OrderNo,
        Status:      4, // Completed
        UserId:      1,
        TotalAmount: event.RefundAmount,
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    // Insert order
    result, err := testCtx.DB.GetOrdersModel().Insert(ctx, order)
    require.NoError(s.T(), err)
    orderId, err := result.LastInsertId()
    require.NoError(s.T(), err)

    // Initialize test data - Payment
    payment := &model.OrderPayments{
        OrderId:       uint64(orderId),
        PaymentNo:     event.PaymentNo,
        PaymentMethod: 1,
        Amount:        event.RefundAmount,
        Status:        1, // Success
        PayTime:       sql.NullTime{Time: time.Now(), Valid: true}, 
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    _, err = testCtx.DB.GetPaymentsModel().Insert(ctx, payment)
    require.NoError(s.T(), err)

    // Initialize test data - Refund
    refund := &model.OrderRefunds{
        OrderId:   uint64(orderId),
        RefundNo:  event.RefundNo,
        Amount:    event.RefundAmount,
        Status:    2, // Processing
        Reason:    event.Reason,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }
    _, err = testCtx.DB.GetRefundsModel().Insert(ctx, refund)
    require.NoError(s.T(), err)

    // Create and bind temporary queue
    queue, err := testCtx.RMQ.GetChannel().QueueDeclare(
        "",    // Empty name
        false, // Non-durable
        true,  // Auto-delete
        true,  // Exclusive
        false,
        nil,
    )
    require.NoError(s.T(), err)

    s.T().Logf("Created temporary queue: %s", queue.Name)

    // Bind to exchange
    err = testCtx.RMQ.GetChannel().QueueBind(
        queue.Name,
        "payment.refund",  
        "payment.events",
        false,
        nil,
    )
    require.NoError(s.T(), err)

    // Publish refund event with correct routing key
    s.T().Logf("Publishing refund event: %+v", event)
    err = testCtx.RMQ.PublishEvent("payment.events", "payment.refund", event)
    require.NoError(s.T(), err)

    // Add verification
    s.T().Log("Waiting for message...")
    msg := helpers.AssertMessageReceived(s.T(), testCtx.RMQ, queue.Name, 5*time.Second)
    s.T().Logf("Received message: %s", string(msg.Body))
    helpers.AssertMessageContent(s.T(), msg, "application/json", "payment.refund")

    // Add retry mechanism for state changes
    s.T().Log("Waiting for order processing...")
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        // Check states
        order, err := testCtx.DB.GetOrdersModel().FindByOrderNo(ctx, event.OrderNo)
        require.NoError(s.T(), err)
        if order.Status == 6 { // Refunded
            break
        }
        s.T().Logf("Current order status: %d, attempt %d/%d", order.Status, i+1, maxRetries)
        time.Sleep(time.Second)
    }

    // Verify final states
    helpers.AssertDatabaseState(s.T(), testCtx.DB, map[string]int64{
        event.OrderNo:   6, // Refunded
        event.PaymentNo: 2, // Refunded
        event.RefundNo:  1, // Success
    })
}
