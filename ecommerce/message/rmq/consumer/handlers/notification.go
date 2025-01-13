package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type NotificationHandler struct {
	messageRpc messageservice.MessageService
	messages   model.MessagesModel
	sends      model.MessageSendsModel
}

func NewNotificationHandler(messageRpc messageservice.MessageService, messages model.MessagesModel, sends model.MessageSendsModel) *NotificationHandler {
	return &NotificationHandler{
		messageRpc: messageRpc,
		messages:   messages,
		sends:      sends,
	}
}

func (h *NotificationHandler) Handle(event *types.MessageEvent) error {
	switch event.Type {
	case types.EventTypeMessageCreated:
		return h.handleMessageCreated(event)
	case types.EventTypeMessageRead:
		return h.handleMessageRead(event)
	default:
		return nil
	}
}

func (h *NotificationHandler) handleMessageCreated(event *types.MessageEvent) error {
	data, ok := event.Data.(*types.MessageCreatedData)
	if !ok {
		return fmt.Errorf("invalid message created data")
	}

	// Start transaction
	err := h.messages.Trans(context.Background(), func(ctx context.Context, session sqlx.Session) error {
		// Insert message
		extraData := sql.NullString{}
		if data.ExtraData != "" {
			extraData.String = data.ExtraData
			extraData.Valid = true
		}

		message := &model.Messages{
			UserId:      uint64(data.UserID),
			Title:       data.Title,
			Content:     data.Content,
			Type:        int64(data.Type),
			SendChannel: int64(data.SendChannel),
			ExtraData:   extraData,
			IsRead:      0,
			CreatedAt:   time.Now(),
		}

		result, err := h.messages.Insert(ctx, message)
		if err != nil {
			return err
		}

		messageId, err := result.LastInsertId()
		if err != nil {
			return err
		}

		// Create send record
		send := &model.MessageSends{
			MessageId:  uint64(messageId),
			UserId:     uint64(data.UserID),
			Channel:    int64(data.SendChannel),
			Status:     types.MessageStatusPending,
			RetryCount: 0,
		}

		_, err = h.sends.Insert(ctx, send)
		return err
	})

	return err
}

func (h *NotificationHandler) handleMessageRead(event *types.MessageEvent) error {
	data, ok := event.Data.(*types.MessageReadData)
	if !ok {
		return fmt.Errorf("invalid message read data")
	}

	message, err := h.messages.FindOne(context.Background(), uint64(data.MessageID))
	if err != nil {
		return err
	}

	message.IsRead = 1
	message.ReadTime = sql.NullTime{
		Time:  data.ReadTime,
		Valid: true,
	}

	return h.messages.Update(context.Background(), message)
}
