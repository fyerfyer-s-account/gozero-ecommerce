package batch

import (
    "encoding/json"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/types"
    "github.com/streadway/amqp"
)

type Sender struct {
    channel  *amqp.Channel
    exchange string
}

func NewSender(channel *amqp.Channel, exchange string) *Sender {
    return &Sender{
        channel:  channel,
        exchange: exchange,
    }
}

func (s *Sender) SendBatch(events []*types.CartEvent) error {
    for _, event := range events {
        body, err := json.Marshal(event)
        if err != nil {
            return err
        }
        
        err = s.channel.Publish(
            s.exchange,
            string(event.Type),
            false,
            false,
            amqp.Publishing{
                ContentType:  "application/json",
                Body:        body,
                MessageId:   event.ID,
                Timestamp:   event.Timestamp,
            },
        )
        if err != nil {
            return err
        }
    }
    return nil
}