// message_test.go
package cases

import (
    "context"
    "testing"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/order/setup"
    "github.com/stretchr/testify/require"
)

func TestMessageDelivery(t *testing.T) {
    // Initialize test context
    _, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    testCtx, err := setup.NewTestContext()
    require.NoError(t, err)
    defer testCtx.Close()

    // Create test queue
    queue, err := testCtx.RMQ.GetChannel().QueueDeclare(
        "test.queue",
        false,
        true,
        false,
        false,
        nil,
    )
    require.NoError(t, err)

    // Bind to exchange
    err = testCtx.RMQ.GetChannel().QueueBind(
        queue.Name,
        "test.key",
        "order.events",
        false,
        nil,
    )
    require.NoError(t, err)

    // Send test message
    testMsg := map[string]string{"test": "message"}
    err = testCtx.RMQ.PublishEvent("order.events", "test.key", testMsg)
    require.NoError(t, err)

    // Verify message received
    msg, err := testCtx.RMQ.ConsumeMessage(queue.Name, 5*time.Second) 
    require.NoError(t, err)
    require.NotNil(t, msg)
}