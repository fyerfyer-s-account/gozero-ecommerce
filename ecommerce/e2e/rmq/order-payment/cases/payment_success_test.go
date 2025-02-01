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

func TestPaymentSuccess(t *testing.T) {
    // Initialize test context
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    testCtx, err := setup.NewTestContext()
    require.NoError(t, err)
    defer testCtx.Close()

    // Load test data
    var event struct {
        OrderNo      string    `json:"order_no"`
        PaymentNo    string    `json:"payment_no"`
        PaymentMethod int64    `json:"payment_method"`
        Amount       float64   `json:"amount"`
        Timestamp    time.Time `json:"timestamp"`
    }
    require.NoError(t, helpers.LoadFixture("payment_success.json", &event))

    // Initialize database state
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

    // Insert payment
    payment := &model.OrderPayments{
        OrderId:       uint64(orderId),
        PaymentNo:     event.PaymentNo,
        PaymentMethod: event.PaymentMethod,
        Amount:        event.Amount,
        Status:        1, // Processing
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    require.NoError(t, helpers.SaveFixture("payment.json", payment))
    _, err = testCtx.DB.GetPaymentsModel().Insert(ctx, payment)
    require.NoError(t, err)

    // Publish payment success event
    err = testCtx.RMQ.PublishEvent("payment.events", "payment.success", event)
    require.NoError(t, err)

    // Verify message received
    msg := helpers.AssertMessageReceived(t, testCtx.RMQ, "order.payment.success", 5*time.Second)
    helpers.AssertMessageContent(t, msg, "application/json", "payment.success")

    // Verify final state
    helpers.AssertDatabaseState(t, testCtx.DB, map[string]int64{
        event.OrderNo:   2, // Paid
        event.PaymentNo: 1, // Success
    })
}