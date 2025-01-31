package handlers

import (
    "context"
    "encoding/json"
    "time"
	"database/sql"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type BatchMessageHandler struct {
    logger             *zerolog.Logger
    messagesModel      model.MessagesModel
    messageSendsModel  model.MessageSendsModel
    templatesModel     model.MessageTemplatesModel
    settingsModel      model.NotificationSettingsModel
}

func NewBatchMessageHandler(
    messagesModel model.MessagesModel,
    messageSendsModel model.MessageSendsModel,
    templatesModel model.MessageTemplatesModel,
    settingsModel model.NotificationSettingsModel,
) *BatchMessageHandler {
    return &BatchMessageHandler{
        logger:             zerolog.GetLogger(),
        messagesModel:      messagesModel,
        messageSendsModel:  messageSendsModel,
        templatesModel:     templatesModel,
        settingsModel:      settingsModel,
    }
}

func (h *BatchMessageHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.MessageBatchEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "batch_id":   event.BatchID,
        "total":      event.Total,
        "completed":  event.Completed,
        "failed":     event.Failed,
        "status":     event.Status,
        "event_type": event.Type,
    }
    h.logger.Info(ctx, "Processing batch message event", fields)

    switch event.Type {
    case types.MessageBatchCreated:
        return h.handleBatchCreated(ctx, event)
    case types.MessageBatchProcessing:
        return h.handleBatchProcessing(ctx, event)
    case types.MessageBatchCompleted:
        return h.handleBatchCompleted(ctx, event)
    case types.MessageBatchFailed:
        return h.handleBatchFailed(ctx, event)
    default:
        return h.updateBatchStatus(ctx, event)
    }
}

func (h *BatchMessageHandler) handleBatchCreated(ctx context.Context, event types.MessageBatchEvent) error {
    // Create batch message record
    message := &model.Messages{
        UserId:      uint64(event.UserID),
        Title:       "Batch Message Created",
        Content:     "New batch message created with ID: " + event.BatchID,
        Type:        3, // Batch message type
        SendChannel: 1, // System notification
        IsRead:      0,
        CreatedAt:   event.Timestamp,
    }

    // Insert message record
    result, err := h.messagesModel.Insert(ctx, message)
    if err != nil {
        return err
    }

    // Get message ID
    messageId, err := result.LastInsertId()
    if err != nil {
        return err
    }

    // Create message send record
    send := &model.MessageSends{
        MessageId:     uint64(messageId),
        UserId:        uint64(event.UserID),
        Channel:       1, // System notification
        Status:        1, // Pending
        RetryCount:    0,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}

func (h *BatchMessageHandler) handleBatchProcessing(ctx context.Context, event types.MessageBatchEvent) error {
    // Update batch processing status
    message := &model.Messages{
        Title:    "Batch Message Processing",
        Content:  "Processing batch message: " + event.BatchID,
        Type:     3,
        UserId:   uint64(event.UserID),
    }

    result, err := h.messagesModel.Insert(ctx, message)
    if err != nil {
        return err
    }

    messageId, err := result.LastInsertId()
    if err != nil {
        return err
    }

    // Create processing status record
    send := &model.MessageSends{
        MessageId:  uint64(messageId),
        UserId:     uint64(event.UserID),
        Channel:    1,
        Status:     2, // Processing
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}

func (h *BatchMessageHandler) handleBatchCompleted(ctx context.Context, event types.MessageBatchEvent) error {
    // Update batch completion status
    message := &model.Messages{
        Title:    "Batch Message Completed",
        Content:  "Completed batch message: " + event.BatchID,
        Type:     3,
        UserId:   uint64(event.UserID),
    }

    result, err := h.messagesModel.Insert(ctx, message)
    if err != nil {
        return err
    }

    messageId, err := result.LastInsertId()
    if err != nil {
        return err
    }

    // Create completion status record
    send := &model.MessageSends{
        MessageId:  uint64(messageId),
        UserId:     uint64(event.UserID),
        Channel:    1,
        Status:     3, // Completed
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
        SendTime:   sql.NullTime{Time: time.Now(), Valid: true},
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}

func (h *BatchMessageHandler) handleBatchFailed(ctx context.Context, event types.MessageBatchEvent) error {
    // Record batch failure
    message := &model.Messages{
        Title:    "Batch Message Failed",
        Content:  "Failed batch message: " + event.BatchID,
        Type:     3,
        UserId:   uint64(event.UserID),
    }

    result, err := h.messagesModel.Insert(ctx, message)
    if err != nil {
        return err
    }

    messageId, err := result.LastInsertId()
    if err != nil {
        return err
    }

    // Create failure record
    send := &model.MessageSends{
        MessageId:  uint64(messageId),
        UserId:     uint64(event.UserID),
        Channel:    1,
        Status:     4, // Failed
        Error:      sql.NullString{String: "Batch processing failed", Valid: true},
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}

func (h *BatchMessageHandler) updateBatchStatus(ctx context.Context, event types.MessageBatchEvent) error {
    message := &model.Messages{
        Title:    "Batch Status Updated",
        Content:  "Updated batch status: " + event.Status,
        Type:     3,
        UserId:   uint64(event.UserID),
    }

    result, err := h.messagesModel.Insert(ctx, message)
    if err != nil {
        return err
    }

    messageId, err := result.LastInsertId()
    if err != nil {
        return err
    }

    // Create status update record
    send := &model.MessageSends{
        MessageId:  uint64(messageId),
        UserId:     uint64(event.UserID),
        Channel:    1,
        Status:     1,
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}