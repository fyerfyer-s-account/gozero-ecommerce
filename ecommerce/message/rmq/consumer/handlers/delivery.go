package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
)

type DeliveryHandler struct {
	sends model.MessageSendsModel
}

func NewDeliveryHandler(sends model.MessageSendsModel) *DeliveryHandler {
	return &DeliveryHandler{
		sends: sends,
	}
}

func (h *DeliveryHandler) Handle(event *types.MessageEvent) error {
	switch event.Type {
	case types.EventTypeMessageSent:
		return h.handleMessageSent(event)
	default:
		return nil
	}
}

func (h *DeliveryHandler) handleMessageSent(event *types.MessageEvent) error {
	data, ok := event.Data.(*types.MessageSentData)
	if !ok {
		return fmt.Errorf("invalid message sent data")
	}

	send, err := h.sends.FindOne(context.Background(), uint64(data.MessageID))
	if err != nil {
		return err
	}

	send.Status = int64(data.Status)
	if data.Error != "" {
		send.Error.String = data.Error
		send.Error.Valid = true
	}

	if data.Status == types.MessageStatusSuccess {
		send.SendTime = sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		}
	} else if data.Status == types.MessageStatusFailed {
		send.RetryCount++

		if send.RetryCount < 3 {
			send.NextRetryTime = sql.NullTime{
				Time:  time.Now().Add(time.Minute * time.Duration(5 * send.RetryCount)),
				Valid: true,
			}
		}
	}

	return h.sends.Update(context.Background(), send)
}
