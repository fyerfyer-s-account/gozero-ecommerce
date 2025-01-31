package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
)

type PaymentSuccessHandler struct {
	logger            *zerolog.Logger
	messagesModel     model.MessagesModel
	messageSendsModel model.MessageSendsModel
	templatesModel    model.MessageTemplatesModel
	settingsModel     model.NotificationSettingsModel
}

func NewPaymentSuccessHandler(
	messagesModel model.MessagesModel,
	messageSendsModel model.MessageSendsModel,
	templatesModel model.MessageTemplatesModel,
	settingsModel model.NotificationSettingsModel,
) *PaymentSuccessHandler {
	return &PaymentSuccessHandler{
		logger:            zerolog.GetLogger(),
		messagesModel:     messagesModel,
		messageSendsModel: messageSendsModel,
		templatesModel:    templatesModel,
		settingsModel:     settingsModel,
	}
}

func (h *PaymentSuccessHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
	var event types.MessagePaymentSuccessEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return err
	}

	fields := map[string]interface{}{
		"order_no":    event.OrderNo,
		"payment_no":  event.PaymentNo,
		"user_id":     event.UserID,
		"amount":      event.Amount,
		"channel":     event.Channel,
		"template_id": event.TemplateID,
		"event_type":  event.Type,
	}
	h.logger.Info(ctx, "Processing payment success notification", fields)

	// Check notification settings
	settings, err := h.settingsModel.FindOneByUserIdTypeChannel(ctx, uint64(event.UserID), 2, 1) // 2: Payment notification, 1: System notification
	if err != nil && err != model.ErrNotFound {
		return err
	}

	if settings != nil && settings.IsEnabled == 0 {
		h.logger.Info(ctx, "Payment notification is disabled for user", fields)
		return nil
	}

	// Get message template
	template, err := h.templatesModel.FindOne(ctx, uint64(event.TemplateID))
	if err != nil {
		return err
	}

	// Create message record
	message := &model.Messages{
		UserId:      uint64(event.UserID),
		Title:       "Payment Success Notification",
		Content:     h.formatPaymentSuccessMessage(template.ContentTemplate, event),
		Type:        2, // Payment message
		SendChannel: 1, // System notification
		IsRead:      0,
		CreatedAt:   event.Timestamp,
	}

	result, err := h.messagesModel.Insert(ctx, message)
	if err != nil {
		return err
	}

	messageId, err := result.LastInsertId()
	if err != nil {
		return err
	}

	// Create message send record
	send := &model.MessageSends{
		MessageId:  uint64(messageId),
		TemplateId: sql.NullInt64{Int64: event.TemplateID, Valid: true},
		UserId:     uint64(event.UserID),
		Channel:    1, // System notification
		Status:     1, // Pending
		RetryCount: 0,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}

	if _, err := h.messageSendsModel.Insert(ctx, send); err != nil {
		return err
	}

	h.logger.Info(ctx, "Payment success notification processed successfully", fields)
	return nil
}

func (h *PaymentSuccessHandler) formatPaymentSuccessMessage(template string, event types.MessagePaymentSuccessEvent) string {
	// For now, return a simple formatted message
	return fmt.Sprintf("Payment successful for order %s\nAmount: %.2f\nPayment No: %s",
		event.OrderNo,
		event.Amount,
		event.PaymentNo,
	)
}
