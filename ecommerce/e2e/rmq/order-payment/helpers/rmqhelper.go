package helpers

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/broker"
	rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
	"github.com/streadway/amqp"
)

type RMQHelper struct {
    broker  broker.Broker
    channel *amqp.Channel
}

// NewRMQHelper creates a new RMQ helper with given configuration
func NewRMQHelper(config *rmqconfig.RabbitMQConfig) (*RMQHelper, error) {
    // Initialize broker
    rmqBroker := broker.NewAMQPBroker(config)
    
    // Connect with timeout
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := rmqBroker.Connect(ctx); err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    // Get channel
    ch, err := rmqBroker.Channel()
    if err != nil {
        rmqBroker.Close()
        return nil, fmt.Errorf("failed to create channel: %w", err)
    }

    helper := &RMQHelper{
        broker:  rmqBroker,
        channel: ch,
    }

    // Setup test queues
    if err := helper.SetupTestQueues(); err != nil {
        helper.Close()
        return nil, err
    }

    return helper, nil
}

// ConsumeMessage consumes a single message from the specified queue with timeout
func (h *RMQHelper) ConsumeMessage(queueName string, timeout time.Duration) (*amqp.Delivery, error) {
    // Create exclusive consumer
    msgs, err := h.channel.Consume(
        queueName,
        "",    // consumer tag
        true,  // auto-ack
        true,  // exclusive
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        return nil, fmt.Errorf("failed to create consumer: %w", err)
    }

    // Wait for message with timeout
    select {
    case msg := <-msgs:
        return &msg, nil
    case <-time.After(timeout):
        return nil, fmt.Errorf("timeout waiting for message on queue %s", queueName)
    }
}

// PublishEvent publishes an event message to the specified exchange
func (h *RMQHelper) PublishEvent(exchange, routingKey string, event interface{}) error {
    // Marshal event to JSON
    body, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %w", err)
    }

    // Publish message
    return h.channel.Publish(
        exchange,
        routingKey,
        false, // mandatory
        false, // immediate
        amqp.Publishing{
            ContentType: "application/json",
            Body:       body,
        },
    )
}

// SetupTestQueues declares required exchanges and queues for testing
func (h *RMQHelper) SetupTestQueues() error {
    // Define test exchanges
    exchanges := []struct {
        name string
        kind string
    }{
        {"order.events", "topic"},
        {"payment.events", "topic"},
    }

    // Declare exchanges
    for _, e := range exchanges {
        err := h.channel.ExchangeDeclare(
            e.name,  // name
            e.kind,  // type
            true,    // durable
            false,   // auto-deleted
            false,   // internal
            false,   // no-wait
            nil,     // arguments
        )
        if err != nil {
            return fmt.Errorf("failed to declare exchange %s: %w", e.name, err)
        }
    }

    // Define test queues
    queues := []struct {
        name       string
        exchange   string
        routingKey string
    }{
        {"order.payment.success", "payment.events", "payment.success"},
        {"order.payment.failed", "payment.events", "payment.failed"},
        {"order.payment.refund", "payment.events", "payment.refund"},
        {"order.status", "order.events", "order.status.*"},
    }

    // Declare queues and bindings
    for _, q := range queues {
        // Declare queue
        _, err := h.channel.QueueDeclare(
            q.name,  // name
            true,    // durable
            false,   // delete when unused
            false,   // exclusive
            false,   // no-wait
            nil,     // arguments
        )
        if err != nil {
            return fmt.Errorf("failed to declare queue %s: %w", q.name, err)
        }

        // Bind queue
        err = h.channel.QueueBind(
            q.name,      // queue name
            q.routingKey, // routing key
            q.exchange,   // exchange
            false,       // no-wait
            nil,         // arguments
        )
        if err != nil {
            return fmt.Errorf("failed to bind queue %s: %w", q.name, err)
        }
    }

    return nil
}

// Close closes the RabbitMQ connection and channel
func (h *RMQHelper) Close() error {
    var errs []error

    if h.channel != nil {
        if err := h.channel.Close(); err != nil {
            errs = append(errs, fmt.Errorf("failed to close channel: %w", err))
        }
    }

    if h.broker != nil {
        if err := h.broker.Close(); err != nil {
            errs = append(errs, fmt.Errorf("failed to close broker: %w", err))
        }
    }

    if len(errs) > 0 {
        return fmt.Errorf("multiple errors during close: %v", errs)
    }
    return nil
}

// GetChannel returns the underlying AMQP channel
func (h *RMQHelper) GetChannel() *amqp.Channel {
    return h.channel
}

// IsConnected checks if the broker connection is still alive
func (h *RMQHelper) IsConnected() bool {
    return h.broker != nil && h.broker.IsConnected()
}