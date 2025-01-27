package broker

import (
    "context"
    "testing"
    "time"
    
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
    "github.com/streadway/amqp"
    "github.com/stretchr/testify/require"
)

func TestReconnection(t *testing.T) {
    ctx := context.Background()
    cfg := config.DefaultConfig()
    broker := NewAMQPBroker(cfg)
    
    // Initial connection
    err := broker.Connect(ctx)
    require.NoError(t, err)
    
    // Force disconnect by closing underlying connection
    broker.conn.Close()
    
    // Wait for reconnection
    time.Sleep(6 * time.Second)
    
    // Verify broker reconnected
    require.True(t, broker.IsConnected())
    
    // Verify can still perform operations
    err = broker.DeclareExchange("test.reconnect", "topic")
    require.NoError(t, err)
}

func TestFullMessageFlow(t *testing.T) {
    ctx := context.Background()
    cfg := config.DefaultConfig()
    
    // Setup publisher
    pub := NewAMQPBroker(cfg)
    err := pub.Connect(ctx)
    require.NoError(t, err)
    defer pub.Close()
    
    // Setup consumer
    cons := NewAMQPBroker(cfg)
    err = cons.Connect(ctx)
    require.NoError(t, err)
    defer cons.Close()
    
    // Create exchange and queue
    exchange := "test.flow.exchange"
    queue := "test.flow.queue"
    key := "test.flow"
    
    err = pub.DeclareExchange(exchange, "topic")
    require.NoError(t, err)
    
    _, err = pub.DeclareQueue(queue)
    require.NoError(t, err)
    
    err = pub.BindQueue(queue, exchange, key)
    require.NoError(t, err)
    
    // Start consuming
    msgs, err := cons.Consume(queue, true)
    require.NoError(t, err)
    
    // Publish multiple messages
    messages := []string{"msg1", "msg2", "msg3"}
    for _, msg := range messages {
        err = pub.Publish(exchange, key, amqp.Publishing{
            ContentType: "text/plain",
            Body:        []byte(msg),
        })
        require.NoError(t, err)
    }
    
    // Verify all messages received in order
    for _, expected := range messages {
        select {
        case msg := <-msgs:
            require.Equal(t, expected, string(msg.Body))
        case <-time.After(5 * time.Second):
            t.Fatal("timeout waiting for message")
        }
    }
}