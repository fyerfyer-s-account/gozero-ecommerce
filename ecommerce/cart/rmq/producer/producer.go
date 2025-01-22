package producer

import (
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/producer/async"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/producer/batch"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/types"
	"github.com/streadway/amqp"
)

type Producer struct {
    conn            *amqp.Connection
    channel         *amqp.Channel
    config          *config.RabbitMQConfig
    batchCollector  *batch.Collector
    asyncDispatcher *async.Dispatcher
    sender          *batch.Sender
}

func NewProducer(config *config.RabbitMQConfig) (*Producer, error) {
    conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
        config.Username,
        config.Password,
        config.Host,
        config.Port,
        config.VHost,
    ))
    if err != nil {
        return nil, err
    }

    channel, err := conn.Channel()
    if err != nil {
        return nil, err
    }

    p := &Producer{
        conn:    conn,
        channel: channel,
        config:  config,
        batchCollector: batch.NewCollector(
            config.Batch.Size,
            time.Duration(config.Batch.FlushInterval)*time.Millisecond,
        ),
        asyncDispatcher: async.NewDispatcher(config.Batch.Workers),
        sender:         batch.NewSender(channel, config.Exchanges.CartEvent.Name),
    }

    return p, nil
}

func (p *Producer) PublishEvent(event *types.CartEvent) error {
    if err := event.Validate(); err != nil {
        return err
    }

    return p.asyncDispatcher.Dispatch(func() error {
        return p.batchCollector.Add(event, p.sender.SendBatch)
    })
}

func (p *Producer) PublishEventSync(event *types.CartEvent) error {
    if err := event.Validate(); err != nil {
        return err
    }

    return p.asyncDispatcher.DispatchSync(func() error {
        return p.sender.SendBatch([]*types.CartEvent{event})
    })
}

func (p *Producer) Close() error {
    p.asyncDispatcher.Wait()
    if err := p.channel.Close(); err != nil {
        return err
    }
    return p.conn.Close()
}