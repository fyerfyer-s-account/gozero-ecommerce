package producer

import (
    "encoding/json"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/types"
    "github.com/streadway/amqp"
    "sync"
)

type Producer struct {
    config  *config.RabbitMQConfig
    conn    *amqp.Connection
    channel *amqp.Channel
    mu      sync.RWMutex
}

func NewProducer(cfg *config.RabbitMQConfig) (*Producer, error) {
    p := &Producer{
        config: cfg,
    }
    if err := p.connect(); err != nil {
        return nil, err
    }
    return p, nil
}

func (p *Producer) connect() error {
    p.mu.Lock()
    defer p.mu.Unlock()

    conn, err := amqp.Dial(p.config.GetDSN())
    if err != nil {
        return err
    }

    ch, err := conn.Channel()
    if err != nil {
        return err
    }

    // Declare exchange
    err = ch.ExchangeDeclare(
        p.config.Exchanges.MarketingEvent.Name,
        p.config.Exchanges.MarketingEvent.Type,
        p.config.Exchanges.MarketingEvent.Durable,
        false, // auto-deleted
        false, // internal
        false, // no-wait
        nil,   // arguments
    )
    if err != nil {
        return err
    }

    p.conn = conn
    p.channel = ch
    return nil
}

func (p *Producer) PublishCouponEvent(event *types.MarketingEvent) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    return p.channel.Publish(
        p.config.Exchanges.MarketingEvent.Name,
        "coupon." + string(event.Type),
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:       data,
        },
    )
}

func (p *Producer) PublishPromotionEvent(event *types.MarketingEvent) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    return p.channel.Publish(
        p.config.Exchanges.MarketingEvent.Name,
        "promotion." + string(event.Type),
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:       data,
        },
    )
}

func (p *Producer) PublishPointsEvent(event *types.MarketingEvent) error {
    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    return p.channel.Publish(
        p.config.Exchanges.MarketingEvent.Name,
        "points." + string(event.Type),
        false,
        false,
        amqp.Publishing{
            ContentType: "application/json",
            Body:       data,
        },
    )
}

func (p *Producer) Close() error {
    p.mu.Lock()
    defer p.mu.Unlock()

    if p.channel != nil {
        p.channel.Close()
    }
    if p.conn != nil {
        p.conn.Close()
    }
    return nil
}