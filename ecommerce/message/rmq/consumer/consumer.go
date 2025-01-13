package consumer

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/types"
    "github.com/streadway/amqp"
)

type MessageHandler interface {
    Handle(event *types.MessageEvent) error
}

type Consumer struct {
    conn         *amqp.Connection
    channel      *amqp.Channel
    config       config.RabbitMQConfig
    handlers     map[string]MessageHandler
    done        chan bool
    notifyClose chan *amqp.Error
}

func NewConsumer(cfg config.RabbitMQConfig) *Consumer {
    return &Consumer{
        config:   cfg,
        handlers: make(map[string]MessageHandler),
        done:    make(chan bool),
    }
}

func (c *Consumer) RegisterHandler(queue string, handler MessageHandler) {
    c.handlers[queue] = handler
}

func (c *Consumer) Start(ctx context.Context) error {
    // Connect to RabbitMQ
    if err := c.connect(); err != nil {
        return err
    }

    // Setup exchanges and queues
    if err := c.setup(); err != nil {
        return err
    }

    // Start consuming from queues
    for queueName, handler := range c.handlers {
        if err := c.consume(ctx, queueName, handler); err != nil {
            return err
        }
    }

    // Monitor connection
    go c.monitor(ctx)

    return nil
}

func (c *Consumer) Stop() {
    c.done <- true
    if c.channel != nil {
        c.channel.Close()
    }
    if c.conn != nil {
        c.conn.Close()
    }
}

func (c *Consumer) connect() error {
    var err error
    c.conn, err = amqp.Dial(c.config.GetDSN())
    if err != nil {
        return fmt.Errorf("failed to connect to RabbitMQ: %v", err)
    }

    c.channel, err = c.conn.Channel()
    if err != nil {
        return fmt.Errorf("failed to open channel: %v", err)
    }

    c.notifyClose = make(chan *amqp.Error)
    c.conn.NotifyClose(c.notifyClose)

    return nil
}

func (c *Consumer) setup() error {
    // Declare exchange
    err := c.channel.ExchangeDeclare(
        c.config.Exchanges.MessageEvent.Name,
        c.config.Exchanges.MessageEvent.Type,
        c.config.Exchanges.MessageEvent.Durable,
        false, // auto-delete
        false, // internal
        false, // no-wait
        nil,   // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare exchange: %v", err)
    }

    // Declare queues
    queues := []config.QueueConfig{
        c.config.Queues.NotificationQueue,
        c.config.Queues.TemplateQueue,
    }

    for _, q := range queues {
        _, err = c.channel.QueueDeclare(
            q.Name,    // name
            q.Durable, // durable
            false,     // delete when unused
            false,     // exclusive
            false,     // no-wait
            nil,       // arguments
        )
        if err != nil {
            return fmt.Errorf("failed to declare queue %s: %v", q.Name, err)
        }

        err = c.channel.QueueBind(
            q.Name,                               // queue name
            q.RoutingKey,                         // routing key
            c.config.Exchanges.MessageEvent.Name, // exchange
            false,
            nil,
        )
        if err != nil {
            return fmt.Errorf("failed to bind queue %s: %v", q.Name, err)
        }
    }

    return nil
}

func (c *Consumer) consume(ctx context.Context, queueName string, handler MessageHandler) error {
    msgs, err := c.channel.Consume(
        queueName, // queue
        "",        // consumer
        false,     // auto-ack
        false,     // exclusive
        false,     // no-local
        false,     // no-wait
        nil,       // args
    )
    if err != nil {
        return fmt.Errorf("failed to register consumer: %v", err)
    }

    go func() {
        for {
            select {
            case <-ctx.Done():
                return
            case msg, ok := <-msgs:
                if !ok {
                    return
                }
                
                // Handle message
                if err := c.handleMessage(msg, handler); err != nil {
                    log.Printf("Error handling message: %v", err)
                    msg.Nack(false, true) // Negative acknowledge, requeue
                } else {
                    msg.Ack(false) // Acknowledge
                }
            }
        }
    }()

    return nil
}

func (c *Consumer) handleMessage(msg amqp.Delivery, handler MessageHandler) error {
    var event types.MessageEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return fmt.Errorf("failed to unmarshal message: %v", err)
    }

    return handler.Handle(&event)
}

func (c *Consumer) monitor(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case <-c.done:
            return
        case err := <-c.notifyClose:
            if err != nil {
                log.Printf("Connection closed: %v", err)
                // Attempt to reconnect
                for {
                    if err := c.connect(); err == nil {
                        if err := c.setup(); err == nil {
                            // Restart consumers
                            for queueName, handler := range c.handlers {
                                if err := c.consume(ctx, queueName, handler); err == nil {
                                    break
                                }
                            }
                        }
                    }
                    time.Sleep(5 * time.Second)
                }
            }
        }
    }
}