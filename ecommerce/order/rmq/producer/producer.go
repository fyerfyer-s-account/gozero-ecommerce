package producer

import (
    "encoding/json"
    "fmt"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
    "github.com/streadway/amqp"
    "time"
)

type Producer struct {
    conn     *amqp.Connection
    channel  *amqp.Channel
    config   *config.RabbitMQConfig
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

    producer := &Producer{
        conn:    conn,
        channel: channel,
        config:  config,
    }

    err = producer.setup()
    if err != nil {
        return nil, err
    }

    return producer, nil
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

func (p *Producer) Close() error {
    if err := p.channel.Close(); err != nil {
        return err
    }
    return p.conn.Close()
}

func (p *Producer) PublishEvent(event *types.OrderEvent) error {
    if event.Timestamp.IsZero() {
        event.Timestamp = time.Now()
    }

    body, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("failed to marshal event: %v", err)
    }

    return p.channel.Publish(
        p.config.Exchanges.OrderEvent.Name, // exchange
        getRoutingKey(event.Type),          // routing key
        false,                              // mandatory
        false,                              // immediate
        amqp.Publishing{
            ContentType:  "application/json",
            Body:        body,
            DeliveryMode: 2, // persistent
            Timestamp:    event.Timestamp,
        },
    )
}

func getRoutingKey(eventType types.EventType) string {
    switch eventType {
    case types.EventTypeOrderCreated:
        return "order.created"
    case types.EventTypeOrderPaid:
        return "order.paid"
    case types.EventTypeOrderCancelled:
        return "order.cancelled"
    case types.EventTypeOrderShipped:
        return "order.shipped"
    case types.EventTypeOrderCompleted:
        return "order.completed"
    default:
        return "order.unknown"
    }
}

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