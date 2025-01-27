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
    conn       *amqp.Connection
    channel    *amqp.Channel
    logger     *zerolog.Logger
    isReady    bool
    mu         sync.RWMutex
    done       chan bool
}

func NewAMQPBroker(cfg *config.RabbitMQConfig) *AMQPBroker {
    return &AMQPBroker{
        config:  cfg,
        logger:  zerolog.GetLogger(),
        done:    make(chan bool),
    }
}

func (b *AMQPBroker) Connect(ctx context.Context) error {
    b.mu.Lock()
    defer b.mu.Unlock()

    if b.isReady {
        return nil
    }

    var err error
    b.conn, err = amqp.DialConfig(b.config.DSN(), amqp.Config{
        Heartbeat: b.config.HeartbeatInterval,
        Locale:    "en_US",
    })
    if err != nil {
        return fmt.Errorf("failed to connect to RabbitMQ: %w", err)
    }

    b.channel, err = b.conn.Channel()
    if err != nil {
        b.conn.Close()
        return fmt.Errorf("failed to open channel: %w", err)
    }

    // Setup QoS
    err = b.channel.Qos(
        b.config.PrefetchCount,
        0,
        b.config.PrefetchGlobal,
    )
    if err != nil {
        b.Close()
        return fmt.Errorf("failed to set QoS: %w", err)
    }

    b.isReady = true
    go b.handleReconnect(ctx)

    return nil
}

func (b *AMQPBroker) handleReconnect(ctx context.Context) {
    for {
        select {
        case <-b.done:
            return
        case <-b.conn.NotifyClose(make(chan *amqp.Error)):
            b.mu.Lock()
            b.isReady = false
            b.mu.Unlock()

            // Reconnect logic
            for {
                select {
                case <-b.done:
                    return
                case <-ctx.Done():
                    return
                default:
                    if err := b.Connect(ctx); err == nil {
                        b.logger.Info(ctx, "Successfully reconnected to RabbitMQ", nil)
                        return
                    }
                    time.Sleep(5 * time.Second)
                }
            }
        }
    }
}

func (b *AMQPBroker) Close() error {
    b.mu.Lock()
    defer b.mu.Unlock()

    if !b.isReady {
        return nil
    }

    close(b.done)

    if b.channel != nil {
        b.channel.Close()
    }

    if b.conn != nil {
        b.conn.Close()
    }

    b.isReady = false
    return nil
}

func (b *AMQPBroker) Publish(exchange, routingKey string, msg amqp.Publishing) error {
    if !b.isReady {
        return fmt.Errorf("not connected to RabbitMQ")
    }

    return b.channel.Publish(
        exchange,
        routingKey,
        false,
        false,
        msg,
    )
}

func (b *AMQPBroker) Consume(queue string, autoAck bool) (<-chan amqp.Delivery, error) {
    if !b.isReady {
        return nil, fmt.Errorf("not connected to RabbitMQ")
    }

    return b.channel.Consume(
        queue,
        "",
        autoAck,
        false,
        false,
        false,
        nil,
    )
}

func (b *AMQPBroker) DeclareExchange(name string, kind string) error {
    if !b.isReady {
        return fmt.Errorf("not connected to RabbitMQ")
    }

    return b.channel.ExchangeDeclare(
        name,
        kind,
        true,  // durable
        false, // auto-deleted
        false, // internal
        false, // no-wait
        nil,   // arguments
    )
}

func (b *AMQPBroker) DeclareQueue(name string) (amqp.Queue, error) {
    if !b.isReady {
        return amqp.Queue{}, fmt.Errorf("not connected to RabbitMQ")
    }

    return b.channel.QueueDeclare(
        name,
        true,  // durable
        false, // auto-delete
        false, // exclusive
        false, // no-wait
        nil,   // arguments
    )
}

func (b *AMQPBroker) BindQueue(queue, exchange, routingKey string) error {
    if !b.isReady {
        return fmt.Errorf("not connected to RabbitMQ")
    }

    return b.channel.QueueBind(
        queue,
        routingKey,
        exchange,
        false,
        nil,
    )
}

func (b *AMQPBroker) Channel() (*amqp.Channel, error) {
    b.mu.RLock()
    defer b.mu.RUnlock()

    if !b.isReady {
        return nil, fmt.Errorf("not connected to RabbitMQ")
    }

    return b.conn.Channel()
}

func (b *AMQPBroker) IsConnected() bool {
    b.mu.RLock()
    defer b.mu.RUnlock()
    return b.isReady
}

func (b *AMQPBroker) Disconnect() error {
    b.mu.Lock()
    defer b.mu.Unlock()

    if !b.isReady {
        return nil
    }

    // Signal reconnect goroutine to stop
    close(b.done)

    // Close channel if open
    if b.channel != nil {
        if err := b.channel.Close(); err != nil {
            b.logger.Error(context.Background(), "Failed to close channel", err, nil)
        }
        b.channel = nil
    }

    // Close connection if open
    if b.conn != nil {
        if err := b.conn.Close(); err != nil {
            b.logger.Error(context.Background(), "Failed to close connection", err, nil)
        }
        b.conn = nil
    }

    b.isReady = false
    return nil
}