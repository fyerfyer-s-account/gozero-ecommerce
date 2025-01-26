package broker

import (
    "context"
    "fmt"
    "sync"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type AMQPBroker struct {
    config     *config.RabbitMQConfig
    connection *amqp.Connection
    channels   sync.Pool
    logger     *zerolog.Logger
    connected  bool
    mu         sync.RWMutex
}

func NewAMQPBroker(config *config.RabbitMQConfig) *AMQPBroker {
    return &AMQPBroker{
        config: config,
        logger: zerolog.GetLogger(),
        channels: sync.Pool{
            New: func() interface{} {
                return nil
            },
        },
    }
}

func (b *AMQPBroker) Connect(ctx context.Context) error {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.IsConnected() {
        return nil
    }

    conn, err := amqp.Dial(b.config.DSN())
    if err != nil {
        return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    b.connection = conn
    b.connected = true

    // Handle connection close
    go func() {
        <-b.connection.NotifyClose(make(chan *amqp.Error))
        b.mu.Lock()
        b.connected = false
        b.mu.Unlock()
        
        // Attempt to reconnect
        for {
            select {
            case <-ctx.Done():
                return
            default:
                time.Sleep(time.Second)
                if err := b.Connect(ctx); err == nil {
                    return
                }
            }
        }
    }()

    return nil
}

func (b *AMQPBroker) Disconnect() error {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.connection != nil {
        if err := b.connection.Close(); err != nil {
            return fmt.Errorf("failed to close connection: %w", err)
        }
        b.connected = false
    }
    return nil
}

func (b *AMQPBroker) Channel() (*amqp.Channel, error) {
    if !b.IsConnected() {
        return nil, fmt.Errorf("not connected to RabbitMQ")
    }

    if ch := b.channels.Get(); ch != nil {
        return ch.(*amqp.Channel), nil
    }

    channel, err := b.connection.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to create channel: %w", err)
    }

    if err := channel.Qos(
        b.config.PrefetchCount,
        0,
        b.config.PrefetchGlobal,
    ); err != nil {
        return nil, fmt.Errorf("failed to set QoS: %w", err)
    }

    return channel, nil
}

func (b *AMQPBroker) DeclareExchange(name, kind string, durable, autoDelete, internal, noWait bool) error {
    channel, err := b.Channel()
    if err != nil {
        return err
    }
    defer channel.Close()

    return channel.ExchangeDeclare(
        name,
        kind,
        durable,
        autoDelete,
        internal,
        noWait,
        nil,
    )
}

func (b *AMQPBroker) DeclareQueue(name string, durable, autoDelete, exclusive, noWait bool) (amqp.Queue, error) {
    channel, err := b.Channel()
    if err != nil {
        return amqp.Queue{}, err
    }
    defer channel.Close()

    return channel.QueueDeclare(
        name,
        durable,
        autoDelete,
        exclusive,
        noWait,
        nil,
    )
}

func (b *AMQPBroker) BindQueue(name, key, exchange string, noWait bool) error {
    channel, err := b.Channel()
    if err != nil {
        return err
    }
    defer channel.Close()

    return channel.QueueBind(
        name,
        key,
        exchange,
        noWait,
        nil,
    )
}

func (b *AMQPBroker) IsConnected() bool {
    b.mu.RLock()
    defer b.mu.RUnlock()
    return b.connected
}