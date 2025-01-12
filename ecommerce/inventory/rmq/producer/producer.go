package producer

import (
    "context"
    "encoding/json"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
    "github.com/streadway/amqp"
    "sync"
)

type Producer struct {
    config   *config.RabbitMQConfig
    conn     *amqp.Connection
    channel  *amqp.Channel
    exchange string
    mu       sync.RWMutex
}

func NewProducer(cfg *config.RabbitMQConfig) (*Producer, error) {
    p := &Producer{
        config:   cfg,
        exchange: cfg.Exchanges.InventoryEvent.Name,
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
        conn.Close()
        return err
    }

    err = ch.ExchangeDeclare(
        p.exchange,
        p.config.Exchanges.InventoryEvent.Type,
        p.config.Exchanges.InventoryEvent.Durable,
        false,
        false,
        false,
        nil,
    )
    if err != nil {
        ch.Close()
        conn.Close()
        return err
    }

    p.conn = conn
    p.channel = ch
    return nil
}

func (p *Producer) PublishStockUpdate(ctx context.Context, data *types.StockUpdateData, userId int64) error {
    event := types.NewInventoryEvent(types.EventTypeStockUpdated, data, userId)
    return p.publishEvent(ctx, event)
}

func (p *Producer) PublishStockAlert(ctx context.Context, data *types.StockAlertData, userId int64) error {
    event := types.NewInventoryEvent(types.EventTypeStockAlert, data, userId)
    return p.publishEvent(ctx, event)
}

func (p *Producer) PublishStockLock(ctx context.Context, data *types.StockLockData, userId int64) error {
    event := types.NewInventoryEvent(types.EventTypeStockLocked, data, userId)
    return p.publishEvent(ctx, event)
}

func (p *Producer) PublishStockUnlock(ctx context.Context, data *types.StockLockData, userId int64) error {
    event := types.NewInventoryEvent(types.EventTypeStockUnlocked, data, userId)
    return p.publishEvent(ctx, event)
}

func (p *Producer) publishEvent(ctx context.Context, event *types.InventoryEvent) error {
    p.mu.RLock()
    defer p.mu.RUnlock()

    data, err := json.Marshal(event)
    if err != nil {
        return err
    }

    return p.channel.Publish(
        p.exchange,
        string(event.Type),
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

    var errs []error
    if p.channel != nil {
        if err := p.channel.Close(); err != nil {
            errs = append(errs, err)
        }
    }
    
    if p.conn != nil {
        if err := p.conn.Close(); err != nil {
            errs = append(errs, err)
        }
    }

    if len(errs) > 0 {
        return errs[0]
    }
    return nil
}