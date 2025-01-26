package broker

import (
    "context"
    "github.com/streadway/amqp"
)

// Broker defines the interface for message brokers
type Broker interface {
    // Connect establishes connection to the broker
    Connect(ctx context.Context) error
    
    // Disconnect closes the connection to the broker
    Disconnect() error
    
    // Channel creates a new channel
    Channel() (*amqp.Channel, error)
    
    // DeclareExchange declares an exchange
    DeclareExchange(name, kind string, durable, autoDelete, internal, noWait bool) error
    
    // DeclareQueue declares a queue
    DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool) (amqp.Queue, error)
    
    // BindQueue binds a queue to an exchange
    BindQueue(name, key, exchange string, noWait bool) error
    
    // IsConnected checks if broker is connected
    IsConnected() bool
}