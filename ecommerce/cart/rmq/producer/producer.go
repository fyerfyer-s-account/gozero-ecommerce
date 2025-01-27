    package producer

    import (
        "context"
        "time"

        "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/producer"
        "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
        "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
        "github.com/streadway/amqp"
    )

    type CartProducer struct {
        logger    *zerolog.Logger
        producer  *producer.Producer
        exchange  string
    }

    func NewCartProducer(ch *amqp.Channel, exchange string) *CartProducer {
        p := producer.NewProducer(ch,
            producer.WithExchange(exchange),
            producer.WithPersistentDelivery(),
            producer.WithAsync(),
            producer.WithBatch(100, time.Second),
        )

        return &CartProducer{
            logger:   zerolog.GetLogger(),
            producer: p,
            exchange: exchange,
        }
    }

    func (p *CartProducer) PublishCartUpdated(ctx context.Context, event *types.CartUpdatedEvent) error {
        return p.producer.Publish(ctx, event)
    }

    func (p *CartProducer) PublishCartCleared(ctx context.Context, event *types.CartClearedEvent) error {
        return p.producer.Publish(ctx, event)
    }

    func (p *CartProducer) PublishCartSelected(ctx context.Context, event *types.CartSelectionEvent) error {
        return p.producer.Publish(ctx, event)
    }

    func (p *CartProducer) Close() error {
        return p.producer.Close()
    }