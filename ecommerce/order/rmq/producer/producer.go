package producer

import (
    "context"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type OrderProducer struct {
    logger    *zerolog.Logger
    producer  *producer.Producer
    exchange  string
}

func NewOrderProducer(ch *amqp.Channel, exchange string) *OrderProducer {
    p := producer.NewProducer(ch,
        producer.WithExchange(exchange),
        producer.WithPersistentDelivery(),
        producer.WithAsync(),
        producer.WithBatch(100, time.Second),
    )

    return &OrderProducer{
        logger:   zerolog.GetLogger(),
        producer: p,
        exchange: exchange,
    }
}

func (p *OrderProducer) PublishStatusChanged(ctx context.Context, event *types.OrderStatusChangedEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *OrderProducer) PublishAlert(ctx context.Context, event *types.OrderAlertEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *OrderProducer) PublishOrderCreated(ctx context.Context, event *types.OrderCreatedEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *OrderProducer) PublishOrderPaid(ctx context.Context, event *types.OrderPaidEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *OrderProducer) PublishOrderShipped(ctx context.Context, event *types.OrderShippedEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *OrderProducer) PublishOrderCompleted(ctx context.Context, event *types.OrderCompletedEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *OrderProducer) Close() error {
    return p.producer.Close()
}