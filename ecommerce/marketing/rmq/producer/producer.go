package producer

import (
    "context"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type MarketingProducer struct {
    logger    *zerolog.Logger
    producer  *producer.Producer
    exchange  string
}

func NewMarketingProducer(ch *amqp.Channel, exchange string) *MarketingProducer {
    p := producer.NewProducer(ch,
        producer.WithExchange(exchange),
        producer.WithPersistentDelivery(),
        producer.WithAsync(),
        producer.WithBatch(100, time.Second),
    )

    return &MarketingProducer{
        logger:   zerolog.GetLogger(),
        producer: p,
        exchange: exchange,
    }
}

func (p *MarketingProducer) PublishCouponEvent(ctx context.Context, event *types.CouponEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *MarketingProducer) PublishPromotionEvent(ctx context.Context, event *types.PromotionEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *MarketingProducer) PublishPointsEvent(ctx context.Context, event *types.PointsEvent) error {
    return p.producer.Publish(ctx, event)
}

func (p *MarketingProducer) Close() error {
    return p.producer.Close()
}