package cases

import (
	"context"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/order-payment/helpers"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/order-payment/setup"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/stretchr/testify/require"
)

func TestPaymentFailedFlow(t *testing.T) {
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
        OrderNo   string    `json:"order_no"`
        PaymentNo string    `json:"payment_no"`
        Amount    float64   `json:"amount"`
        Reason    string    `json:"reason"`
        ErrorCode string    `json:"error_code"`
        Timestamp time.Time `json:"timestamp"`
    }
    require.NoError(t, helpers.LoadFixture("payment_failed.json", &event))

    // Initialize test data - Order
    order := &model.Orders{
        OrderNo:     event.OrderNo,
        Status:      1, // Pending Payment
        UserId:      1,
        TotalAmount: event.Amount,
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
        PaymentMethod: 1, // Assuming default payment method
        Amount:        event.Amount,
        Status:        1, // Processing
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    _, err = testCtx.DB.GetPaymentsModel().Insert(ctx, payment)
    require.NoError(t, err)

    // Publish payment failed event
    err = testCtx.RMQ.PublishEvent("payment.events", "payment.failed", event)
    require.NoError(t, err)

    // Verify message received
    msg := helpers.AssertMessageReceived(t, testCtx.RMQ, "order.payment.failed", 5*time.Second)
    helpers.AssertMessageContent(t, msg, "application/json", "payment.failed")

    // Verify final state
    helpers.AssertDatabaseState(t, testCtx.DB, map[string]int64{
        event.OrderNo:   5, // Payment Failed
        event.PaymentNo: 0, // Failed
    })
}