package producer

import (
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/producer/async"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/producer/batch"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
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
    conn, err := amqp.Dial(config.GetDSN())
    if err != nil {
        return nil, fmt.Errorf("failed to connect to RabbitMQ: %v", err)
    }

    channel, err := conn.Channel()
    if err != nil {
        return nil, fmt.Errorf("failed to open channel: %v", err)
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
        sender:         batch.NewSender(channel, config.Exchanges.OrderEvent.Name),
    }

    if err := p.setup(); err != nil {
        return nil, err
    }

    return p, nil
}

func (p *Producer) PublishEvent(event *types.OrderEvent) error {
    if err := event.Validate(); err != nil {
        return err
    }

    return p.asyncDispatcher.Dispatch(func() error {
        return p.batchCollector.Add(event, p.sender.SendBatch)
    })
}

func (p *Producer) PublishEventSync(event *types.OrderEvent) error {
    if err := event.Validate(); err != nil {
        return err
    }

    return p.asyncDispatcher.DispatchSync(func() error {
        return p.sender.SendBatch([]*types.OrderEvent{event})
    })
}

func (p *Producer) Close() error {
    p.asyncDispatcher.Wait()
    if err := p.channel.Close(); err != nil {
        return err
    }
    return p.conn.Close()
}

// func getRoutingKey(eventType types.EventType) string {
//     switch eventType {
//     case types.EventTypeOrderCreated:
//         return "order.created"
//     case types.EventTypeOrderPaid:
//         return "order.paid"
//     case types.EventTypeOrderCancelled:
//         return "order.cancelled"
//     case types.EventTypeOrderShipped:
//         return "order.shipped"
//     case types.EventTypeOrderCompleted:
//         return "order.completed"
//     default:
//         return "order.unknown"
//     }
// }

// Helper methods for creating different types of events
func CreateOrderEvent(id string, eventType types.EventType, data interface{}, metadata types.Metadata) *types.OrderEvent {
    return &types.OrderEvent{
        ID:        id,
        Type:      eventType,
        Timestamp: time.Now(),
        Data:      data,
        Metadata:  metadata,
    }
}

func (p *Producer) setup() error {
    return p.channel.ExchangeDeclare(
        p.config.Exchanges.OrderEvent.Name,
        p.config.Exchanges.OrderEvent.Type,
        p.config.Exchanges.OrderEvent.Durable,
        false, // auto-deleted
        false, // internal
        false, // no-wait
        nil,   // arguments
    )
}