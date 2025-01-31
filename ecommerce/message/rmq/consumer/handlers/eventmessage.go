package handlers

import (
    "context"
    "database/sql"
    "encoding/json"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type EventMessageHandler struct {
    logger             *zerolog.Logger
    messagesModel      model.MessagesModel
    messageSendsModel  model.MessageSendsModel
    templatesModel     model.MessageTemplatesModel
    settingsModel      model.NotificationSettingsModel
}

func NewEventMessageHandler(
    messagesModel model.MessagesModel,
    messageSendsModel model.MessageSendsModel,
    templatesModel model.MessageTemplatesModel,
    settingsModel model.NotificationSettingsModel,
) *EventMessageHandler {
    return &EventMessageHandler{
        logger:             zerolog.GetLogger(),
        messagesModel:      messagesModel,
        messageSendsModel:  messageSendsModel,
        templatesModel:     templatesModel,
        settingsModel:      settingsModel,
    }
}

func (h *EventMessageHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.MessageEventSentEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "message_id": event.MessageID,
        "user_id":    event.UserID,
        "channel":    event.Channel,
        "status":     event.Status,
        "event_type": event.Type,
    }
    h.logger.Info(ctx, "Processing message event", fields)

    switch event.Type {
    case types.MessageEventSent:
        return h.handleMessageSent(ctx, event)
    case types.MessageEventReceived:
        return h.handleMessageReceived(ctx, event)
    case types.MessageEventFailed:
        return h.handleMessageFailed(ctx, event)
    default:
        return h.updateMessageStatus(ctx, event)
    }
}

func (h *EventMessageHandler) handleMessageSent(ctx context.Context, event types.MessageEventSentEvent) error {
    // Create message record
    message := &model.Messages{
        UserId:      uint64(event.UserID),
        Title:       "Message Sent",
        Content:     event.Content,
        Type:        1, // System message
        SendChannel: h.getChannelType(event.Channel),
        IsRead:      0,
        CreatedAt:   event.Timestamp,
    }

    // Insert message record
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
        MessageId:     uint64(messageId),
        UserId:        uint64(event.UserID),
        Channel:       h.getChannelType(event.Channel),
        Status:        3, // Sent successfully
        SendTime:      sql.NullTime{Time: time.Now(), Valid: true},
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}

func (h *EventMessageHandler) handleMessageReceived(ctx context.Context, event types.MessageEventSentEvent) error {
    // Update message read status
    message := &model.Messages{
        UserId:    uint64(event.UserID),
        IsRead:    1,
        ReadTime:  sql.NullTime{Time: time.Now(), Valid: true},
        CreatedAt: event.Timestamp,
    }

    result, err := h.messagesModel.Insert(ctx, message)
    if err != nil {
        return err
    }

    messageId, err := result.LastInsertId()
    if err != nil {
        return err
    }

    // Update message send status
    send := &model.MessageSends{
        MessageId:  uint64(messageId),
        UserId:     uint64(event.UserID),
        Channel:    h.getChannelType(event.Channel),
        Status:     3, // Received
        SendTime:   sql.NullTime{Time: time.Now(), Valid: true},
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}

func (h *EventMessageHandler) handleMessageFailed(ctx context.Context, event types.MessageEventSentEvent) error {
    // Create failed message record
    message := &model.Messages{
        UserId:      uint64(event.UserID),
        Title:       "Message Failed",
        Content:     event.Content,
        Type:        1, // System message
        SendChannel: h.getChannelType(event.Channel),
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

    // Create failed send record
    send := &model.MessageSends{
        MessageId:  uint64(messageId),
        UserId:     uint64(event.UserID),
        Channel:    h.getChannelType(event.Channel),
        Status:     4, // Failed
        Error:      sql.NullString{String: "Message delivery failed", Valid: true},
        RetryCount: 0,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}

func (h *EventMessageHandler) updateMessageStatus(ctx context.Context, event types.MessageEventSentEvent) error {
    // Update message status based on event
    message := &model.Messages{
        UserId:      uint64(event.UserID),
        Content:     event.Content,
        Type:        1, // System message
        SendChannel: h.getChannelType(event.Channel),
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

    // Update send status
    send := &model.MessageSends{
        MessageId:  uint64(messageId),
        UserId:     uint64(event.UserID),
        Channel:    h.getChannelType(event.Channel),
        Status:     h.getStatusFromEvent(event.Status),
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}

// Helper functions
func (h *EventMessageHandler) getChannelType(channel string) int64 {
    switch channel {
    case "sms":
        return 2
    case "email":
        return 3
    case "push":
        return 4
    default:
        return 1 // Default to system notification
    }
}

func (h *EventMessageHandler) getStatusFromEvent(status string) int64 {
    switch status {
    case "pending":
        return 1
    case "processing":
        return 2
    case "success":
        return 3
    case "failed":
        return 4
    default:
        return 1
    }
}