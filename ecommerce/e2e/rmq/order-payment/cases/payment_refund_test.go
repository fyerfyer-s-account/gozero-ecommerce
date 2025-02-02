package cases

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/order-payment/helpers"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/order-payment/setup"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/stretchr/testify/require"
)

func TestPaymentRefundFlow(t *testing.T) {
    // Initialize test context with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    testCtx, err := setup.NewTestContext()
    require.NoError(t, err)
    defer testCtx.Close()

    // Clean test data before start
    require.NoError(t, testCtx.DB.CleanTestData(ctx))

    // Load test event fixture
    var event struct {
        OrderNo      string    `json:"order_no"`
        PaymentNo    string    `json:"payment_no"`
        RefundNo     string    `json:"refund_no"`
        RefundAmount float64   `json:"refund_amount"`
        Reason       string    `json:"reason"`
        Timestamp    time.Time `json:"timestamp"`
    }
    require.NoError(t, helpers.LoadFixture("payment_refund.json", &event))

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
    require.NoError(t, err)
    orderId, err := result.LastInsertId()
    require.NoError(t, err)

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
    require.NoError(t, err)

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
    require.NoError(t, err)

    // Create and bind temporary queue
    queue, err := testCtx.RMQ.GetChannel().QueueDeclare(
        "",    // Empty name
        false, // Non-durable
        true,  // Auto-delete
        true,  // Exclusive
        false,
        nil,
    )
    require.NoError(t, err)

    t.Logf("Created temporary queue: %s", queue.Name)

    // Bind to exchange
    err = testCtx.RMQ.GetChannel().QueueBind(
        queue.Name,
        "payment.refund",  
        "payment.events",
        false,
        nil,
    )
    require.NoError(t, err)

    // Publish refund event with correct routing key
    t.Logf("Publishing refund event: %+v", event)
    err = testCtx.RMQ.PublishEvent("payment.events", "payment.refund", event)
    require.NoError(t, err)

    // Add verification
    t.Log("Waiting for message...")
    msg := helpers.AssertMessageReceived(t, testCtx.RMQ, queue.Name, 5*time.Second)
    t.Logf("Received message: %s", string(msg.Body))
    helpers.AssertMessageContent(t, msg, "application/json", "payment.refund")

    // Add retry mechanism for state changes
    t.Log("Waiting for order processing...")
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        // Check states
        order, err := testCtx.DB.GetOrdersModel().FindByOrderNo(ctx, event.OrderNo)
        require.NoError(t, err)
        if order.Status == 6 { // Refunded
            break
        }
        t.Logf("Current order status: %d, attempt %d/%d", order.Status, i+1, maxRetries)
        time.Sleep(time.Second)
    }

    // Verify final states
    helpers.AssertDatabaseState(t, testCtx.DB, map[string]int64{
        event.OrderNo:   6, // Refunded
        event.PaymentNo: 2, // Refunded
        event.RefundNo:  1, // Success
    })
}