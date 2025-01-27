package broker

import (
    "context"
    "testing"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
    "github.com/streadway/amqp"
    "github.com/stretchr/testify/require"
)

func TestManualPublishMessage(t *testing.T) {
    // Setup
    cfg := config.DefaultConfig()
    b := NewAMQPBroker(cfg)
    
    err := b.Connect(context.Background())
    require.NoError(t, err, "Should connect to RabbitMQ")
    defer b.Close()
    
    require.True(t, b.IsConnected(), "Should be connected to RabbitMQ")

    // 1. Declare a test queue
    queue, err := b.DeclareQueue("test.manual.queue")
    require.NoError(t, err, "Should declare queue")

    // 2. Bind queue to amq.topic exchange
    err = b.BindQueue(queue.Name, "amq.topic", "test.key")
    require.NoError(t, err, "Should bind queue")

    // 3. Create consumer from the queue
    msgs, err := b.Consume(queue.Name, true)
    require.NoError(t, err, "Should create consumer")

    // 4. Publish test message
    testMsg := "Hello RabbitMQ"
    err = b.Publish("amq.topic", "test.key", amqp.Publishing{
        ContentType: "text/plain",
        Body:       []byte(testMsg),
    })
    require.NoError(t, err, "Should publish message")

    // 5. Verify message received
    select {
    case msg := <-msgs:
        require.Equal(t, testMsg, string(msg.Body), "Should receive correct message")
    case <-time.After(5 * time.Second):
        t.Fatal("Timeout waiting for message")
    }
}