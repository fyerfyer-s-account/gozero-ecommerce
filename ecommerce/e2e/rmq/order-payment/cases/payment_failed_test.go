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

	// Create a temporary test queue
	queue, err := testCtx.RMQ.GetChannel().QueueDeclare(
		"",    // Empty name, automatically generated
		false, // Non-durable
		true,  // Auto-delete
		true,  // Exclusive
		false,
		nil,
	)
	require.NoError(t, err)

	t.Logf("Created temporary queue: %s", queue.Name)	

	// Bind to the payment.events exchange
	err = testCtx.RMQ.GetChannel().QueueBind(
		queue.Name,       // Use generated queue name
		"payment.failed",
		"payment.events",
		false,
		nil,
	)
	require.NoError(t, err)

	// Publish payment failed event
	t.Logf("Publishing payment failed event: %+v", event)
    err = testCtx.RMQ.PublishEvent("payment.events", "payment.failed", event)
    require.NoError(t, err)

	// Verify message received using the generated queue name
	t.Log("Waiting for message...")
    msg := helpers.AssertMessageReceived(t, testCtx.RMQ, queue.Name, 5*time.Second)
    t.Logf("Received message: %s", string(msg.Body))
	helpers.AssertMessageContent(t, msg, "application/json", "payment.failed")

	time.Sleep(2 * time.Second)

	t.Log("Waiting for order processing...")
    maxRetries := 5
    for i := 0; i < maxRetries; i++ {
        // Check order status
        order, err := testCtx.DB.GetOrdersModel().FindByOrderNo(ctx, event.OrderNo)
        require.NoError(t, err)
        if order.Status == 5 { // Payment Failed
            break
        }
        t.Logf("Current order status: %d, attempt %d/%d", order.Status, i+1, maxRetries)
        time.Sleep(time.Second)
    }

    // Check payment status
    payment, err = testCtx.DB.GetPaymentsModel().FindOneByPaymentNo(ctx, event.PaymentNo)
    require.NoError(t, err)
    t.Logf("Final payment status: %d", payment.Status)

	// Verify final state
	helpers.AssertDatabaseState(t, testCtx.DB, map[string]int64{
		event.OrderNo:   5, // Payment Failed
		event.PaymentNo: 0, // Failed
	})
}
