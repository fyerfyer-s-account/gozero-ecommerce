package consumer

import (
    "encoding/json"
    "fmt"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
    "github.com/streadway/amqp"
    "log"
)

type EventHandler interface {
    Handle(event *types.OrderEvent) error
}

type Consumer struct {
    conn      *amqp.Connection
    channel   *amqp.Channel
    config    *config.RabbitMQConfig
    handlers  map[types.EventType][]EventHandler
}

func NewConsumer(config *config.RabbitMQConfig) (*Consumer, error) {
    conn, err := amqp.Dial(config.GetDSN())
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
    }

    channel, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to open channel: %v", err)
    }

    return &Consumer{
        conn:     conn,
        channel:  channel,
        config:   config,
        handlers: make(map[types.EventType][]EventHandler),
    }, nil
}

func (c *Consumer) Subscribe(eventType types.EventType, handler EventHandler) {
    c.handlers[eventType] = append(c.handlers[eventType], handler)
}

func (c *Consumer) Start() error {
    err := c.setupExchangesAndQueues()
    if err != nil {
        return err
    }

    msgs, err := c.channel.Consume(
        c.config.Queues.OrderStatus.Name,
        "",    // consumer
        true,  // auto-ack
        false, // exclusive
        false, // no-local
        false, // no-wait
        nil,   // args
    )
    if err != nil {
        return fmt.Errorf("failed to register a consumer: %v", err)
    }

    go c.handleMessages(msgs)
    return nil
}

func (c *Consumer) setupExchangesAndQueues() error {
    // Declare exchange
    err := c.channel.ExchangeDeclare(
        c.config.Exchanges.OrderEvent.Name,
        c.config.Exchanges.OrderEvent.Type,
        c.config.Exchanges.OrderEvent.Durable,
        false, // auto-deleted
        false, // internal
        false, // no-wait
        nil,   // arguments
    )
    if err != nil {
        return fmt.Errorf("failed to declare exchange: %v", err)
    }

    // Declare queues and bind them
    for _, queueConfig := range []config.QueueConfig{
        c.config.Queues.OrderStatus,
        c.config.Queues.OrderAlert,
    } {
        _, err = c.channel.QueueDeclare(
            queueConfig.Name,
            queueConfig.Durable,
            false, // auto-delete
            false, // exclusive
            false, // no-wait
            nil,   // arguments
        )
        if err != nil {
            return fmt.Errorf("failed to declare queue %s: %v", queueConfig.Name, err)
        }

        err = c.channel.QueueBind(
            queueConfig.Name,
            queueConfig.RoutingKey,
            c.config.Exchanges.OrderEvent.Name,
            false,
            nil,
        )
        if err != nil {
            return fmt.Errorf("failed to bind queue %s: %v", queueConfig.Name, err)
        }
    }

    return nil
}

func (c *Consumer) handleMessages(msgs <-chan amqp.Delivery) {
    for msg := range msgs {
        var event types.OrderEvent
        err := json.Unmarshal(msg.Body, &event)
        if err != nil {
            log.Printf("Error unmarshaling message: %v", err)
            continue
        }

        handlers, exists := c.handlers[event.Type]
        if !exists {
            continue
        }

        for _, handler := range handlers {
            if err := handler.Handle(&event); err != nil {
                log.Printf("Error handling event %s: %v", event.Type, err)
            }
        }
    }
}

func (c *Consumer) Close() error {
    if err := c.channel.Close(); err != nil {
        return err
    }
    return c.conn.Close()
}