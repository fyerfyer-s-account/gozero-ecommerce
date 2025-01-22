package batch

import (
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
    "github.com/streadway/amqp"
    "encoding/json"
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

func (s *Sender) SendBatch(events []*types.InventoryEvent) error {
    if len(events) == 0 {
        return nil
    }

    for _, event := range events {
        body, err := json.Marshal(event)
        if err != nil {
            return err
        }

        msg := amqp.Publishing{
            ContentType: "application/json",
            Body:       body,
            MessageId:  event.ID,
            Timestamp:  event.Timestamp,
        }

        err = s.channel.Publish(
            s.exchange,
            string(event.Type),
            false,
            false,
            msg,
        )
        if err != nil {
            return err
        }
    }
    return nil
}