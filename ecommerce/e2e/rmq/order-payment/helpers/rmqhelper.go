package helpers

import (
    "context"
    "fmt"
    "time"

    "github.com/streadway/amqp"
)

type RMQHelper struct {
    conn    *amqp.Connection
    channel *amqp.Channel
}

func NewRMQHelper() (*RMQHelper, error) {
    // Connect to local RabbitMQ instance
    conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, fmt.Errorf("failed to open channel: %w", err)
    }

    return &RMQHelper{
        conn:    conn,
        channel: ch,
    }, nil
}

// ConsumeMessage consumes a single message from the specified queue with timeout
func (h *RMQHelper) ConsumeMessage(queueName string, timeout time.Duration) (*amqp.Delivery, error) {
    msgs, err := h.channel.Consume(
        queueName,
        "",    // consumer
        false, // auto-ack
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        return nil, fmt.Errorf("failed to consume from queue: %w", err)
    }

    ctx, cancel := context.WithTimeout(context.Background(), timeout)
    defer cancel()

    select {
    case msg := <-msgs:
        return &msg, nil
    case <-ctx.Done():
        return nil, fmt.Errorf("timeout waiting for message from queue %s", queueName)
    }
}

// PublishMessage publishes a message to the specified exchange
func (h *RMQHelper) PublishMessage(exchange, routingKey string, message []byte) error {
    return h.channel.Publish(
        exchange,
        routingKey,
        false, // mandatory
        false, // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:       message,
        },
    )
}

// Close closes the RabbitMQ connection and channel
func (h *RMQHelper) Close() error {
    if h.channel != nil {
        if err := h.channel.Close(); err != nil {
            return fmt.Errorf("failed to close channel: %w", err)
        }
    }
    if h.conn != nil {
        if err := h.conn.Close(); err != nil {
            return fmt.Errorf("failed to close connection: %w", err)
        }
    }
    return nil
}