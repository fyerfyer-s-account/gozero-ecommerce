package consumer

import (
    "encoding/json"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/consumer/handlers"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/types"
    "github.com/streadway/amqp"
    "log"
    "sync"
)

type Consumer struct {
    config   *config.RabbitMQConfig
    conn     *amqp.Connection
    channel  *amqp.Channel
    handlers map[types.EventType]func(*types.MarketingEvent) error
    mu       sync.RWMutex
}

func NewConsumer(cfg *config.RabbitMQConfig) (*Consumer, error) {
    c := &Consumer{
        config:   cfg,
        handlers: make(map[types.EventType]func(*types.MarketingEvent) error),
    }

    if err := c.connect(); err != nil {
        return nil, err
    }

    c.registerHandlers()
    return c, nil
}

func (c *Consumer) connect() error {
    c.mu.Lock()
    defer c.mu.Unlock()

    conn, err := amqp.Dial(c.config.GetDSN())
    if err != nil {
        return err
    }

    ch, err := conn.Channel()
    if err != nil {
        return err
    }

    c.conn = conn
    c.channel = ch
    return nil
}

func (c *Consumer) registerHandlers() {
    couponHandler := handlers.NewCouponHandler()
    promotionHandler := handlers.NewPromotionHandler()
    pointsHandler := handlers.NewPointsHandler()

    c.handlers[types.EventTypeCouponReceived] = couponHandler.HandleCouponReceived
    c.handlers[types.EventTypeCouponUsed] = couponHandler.HandleCouponUsed
    c.handlers[types.EventTypePromotionStarted] = promotionHandler.HandlePromotionStatus
    c.handlers[types.EventTypePromotionEnded] = promotionHandler.HandlePromotionStatus
    c.handlers[types.EventTypePointsAdded] = pointsHandler.HandlePointsTransaction
    c.handlers[types.EventTypePointsUsed] = pointsHandler.HandlePointsTransaction
}

func (c *Consumer) Start() error {
    for _, q := range []struct {
        name       string
        routingKey string
    }{
        {c.config.Queues.CouponEvent.Name, c.config.Queues.CouponEvent.RoutingKey},
        {c.config.Queues.PromotionEvent.Name, c.config.Queues.PromotionEvent.RoutingKey},
        {c.config.Queues.PointsEvent.Name, c.config.Queues.PointsEvent.RoutingKey},
    } {
        queue, err := c.channel.QueueDeclare(
            q.name,
            true,
            false,
            false,
            false,
            nil,
        )
        if err != nil {
            return err
        }

        err = c.channel.QueueBind(
            queue.Name,
            q.routingKey,
            c.config.Exchanges.MarketingEvent.Name,
            false,
            nil,
        )
        if err != nil {
            return err
        }

        msgs, err := c.channel.Consume(
            queue.Name,
            "",
            true,
            false,
            false,
            false,
            nil,
        )
        if err != nil {
            return err
        }

        go c.handle(msgs)
    }

    return nil
}

func (c *Consumer) handle(msgs <-chan amqp.Delivery) {
    for msg := range msgs {
        var event types.MarketingEvent
        if err := json.Unmarshal(msg.Body, &event); err != nil {
            log.Printf("Error unmarshaling message: %v", err)
            continue
        }

        if handler, ok := c.handlers[event.Type]; ok {
            if err := handler(&event); err != nil {
                log.Printf("Error handling message: %v", err)
            }
        }
    }
}

func (c *Consumer) Close() error {
    c.mu.Lock()
    defer c.mu.Unlock()

    if c.channel != nil {
        c.channel.Close()
    }
    if c.conn != nil {
        c.conn.Close()
    }
    return nil
}