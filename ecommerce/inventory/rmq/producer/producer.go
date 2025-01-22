package producer

import (
    "context"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/producer/async"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/producer/batch"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
    "github.com/streadway/amqp"
    "time"
)

type Producer struct {
    conn            *amqp.Connection
    channel         *amqp.Channel
    config          *config.RabbitMQConfig
    batchCollector  *batch.Collector
    asyncDispatcher *async.Dispatcher
    sender          *batch.Sender
}

func NewProducer(cfg *config.RabbitMQConfig) (*Producer, error) {
    conn, err := amqp.Dial(cfg.GetDSN())
    if err != nil {
        return nil, err
    }

    ch, err := conn.Channel()
    if err != nil {
        conn.Close()
        return nil, err
    }

    p := &Producer{
        conn:    conn,
        channel: ch,
        config:  cfg,
        batchCollector: batch.NewCollector(
            cfg.Batch.Size,
            time.Duration(cfg.Batch.FlushInterval)*time.Millisecond,
        ),
        asyncDispatcher: async.NewDispatcher(cfg.Batch.Workers),
        sender:         batch.NewSender(ch, cfg.Exchanges.InventoryEvent.Name),
    }

    if err := p.setup(); err != nil {
        return nil, err
    }

    return p, nil
}

func (p *Producer) setup() error {
    return p.channel.ExchangeDeclare(
        p.config.Exchanges.InventoryEvent.Name,
        p.config.Exchanges.InventoryEvent.Type,
        p.config.Exchanges.InventoryEvent.Durable,
        false,
        false,
        false,
        nil,
    )
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
    return p.asyncDispatcher.Dispatch(func() error {
        return p.batchCollector.Add(event, p.sender.SendBatch)
    })
}

func (p *Producer) publishEventSync(ctx context.Context, event *types.InventoryEvent) error {
    return p.asyncDispatcher.DispatchSync(func() error {
        return p.sender.SendBatch([]*types.InventoryEvent{event})
    })
}

func (p *Producer) Close() error {
    p.asyncDispatcher.Wait()
    if err := p.channel.Close(); err != nil {
        return err
    }
    return p.conn.Close()
}