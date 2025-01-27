package broker

import (
    "context"
    "github.com/streadway/amqp"
)

// Broker defines the interface for a message broker
type Broker interface {
    // Connect connects to the message queue
    Connect(ctx context.Context) error
    // Close closes the connection
    Close() error
    // Publish publishes a message
    Publish(exchange, routingKey string, msg amqp.Publishing) error
    // Consume consumes messages from a queue
    Consume(queue string, autoAck bool) (<-chan amqp.Delivery, error)
    // DeclareExchange declares an exchange
    DeclareExchange(name string, kind string) error
    // DeclareQueue declares a queue
    DeclareQueue(name string) (amqp.Queue, error)
    // BindQueue binds a queue to an exchange
    BindQueue(queue, exchange, routingKey string) error
    // IsConnected checks the connection status
    IsConnected() bool

    Channel() (*amqp.Channel, error)
}
