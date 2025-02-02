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

    // Clean test data
    require.NoError(t, testCtx.DB.CleanTestData(ctx))

    // Load test data
    var event struct {
        OrderNo       string    `json:"order_no"`
        PaymentNo     string    `json:"payment_no"`
        PaymentMethod int64     `json:"payment_method"`
        Amount        float64   `json:"amount"`
        Timestamp     time.Time `json:"timestamp"`
    }
    require.NoError(t, helpers.LoadFixture("payment_success.json", &event))

    // Initialize database state - Order
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

    // Initialize database state - Payment
    payment := &model.OrderPayments{
        OrderId:       uint64(orderId),
        PaymentNo:     event.PaymentNo,
        PaymentMethod: event.PaymentMethod,
        Amount:        event.Amount,
        Status:        1, // Processing
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    _, err = testCtx.DB.GetPaymentsModel().Insert(ctx, payment)
    require.NoError(t, err)

    // Create temporary queue
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

    // Bind queue to exchange
    err = testCtx.RMQ.GetChannel().QueueBind(
        queue.Name,
        "payment.success",
        "payment.events",
        false,
        nil,
    )
    require.NoError(t, err)

    // Publish payment success event
    t.Logf("Publishing payment success event: %+v", event)
    err = testCtx.RMQ.PublishEvent("payment.events", "payment.success", event)
    require.NoError(t, err)

    // Verify message received
    t.Log("Waiting for message...")
    msg := helpers.AssertMessageReceived(t, testCtx.RMQ, queue.Name, 5*time.Second)
    t.Logf("Received message: %s", string(msg.Body))
    helpers.AssertMessageContent(t, msg, "application/json", "payment.success")

    // Add retry mechanism for state changes
    t.Log("Waiting for order processing...")
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        // Check states
        order, err := testCtx.DB.GetOrdersModel().FindByOrderNo(ctx, event.OrderNo)
        require.NoError(t, err)
        if order.Status == 2 { // Paid
            break
        }
        t.Logf("Current order status: %d, attempt %d/%d", order.Status, i+1, maxRetries)
        time.Sleep(time.Second)
    }

    // Verify final states
    helpers.AssertDatabaseState(t, testCtx.DB, map[string]int64{
        event.OrderNo:   2, // Paid
        event.PaymentNo: 1, // Success
    })
}