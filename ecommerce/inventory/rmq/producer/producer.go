package producer

import (
    "context"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type InventoryProducer struct {
    logger    *zerolog.Logger
    producer  *producer.Producer
    exchange  string
}

func NewInventoryProducer(ch *amqp.Channel, exchange string) *InventoryProducer {
    p := producer.NewProducer(ch,
        producer.WithExchange(exchange),
        producer.WithPersistentDelivery(),
        producer.WithAsync(),
        producer.WithBatch(100, time.Second),
    )

    return &InventoryProducer{
        logger:   zerolog.GetLogger(),
        producer: p,
        exchange: exchange,
    }
}

func (p *InventoryProducer) PublishStockUpdated(ctx context.Context, event *types.StockUpdatedEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *InventoryProducer) PublishStockOutOfStock(ctx context.Context, event *types.StockOutOfStockEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *InventoryProducer) PublishStockLowStock(ctx context.Context, event *types.StockLowStockEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *InventoryProducer) PublishStockLocked(ctx context.Context, event *types.StockLockedEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *InventoryProducer) PublishStockUnlocked(ctx context.Context, event *types.StockUnlockedEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *InventoryProducer) PublishStockDeducted(ctx context.Context, event *types.StockDeductedEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *InventoryProducer) Close() error {
    return p.producer.Close()
}