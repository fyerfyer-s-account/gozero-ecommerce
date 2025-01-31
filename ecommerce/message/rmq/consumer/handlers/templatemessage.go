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

type TemplateMessageHandler struct {
    logger             *zerolog.Logger
    messagesModel      model.MessagesModel
    messageSendsModel  model.MessageSendsModel
    templatesModel     model.MessageTemplatesModel
    settingsModel      model.NotificationSettingsModel
}

func NewTemplateMessageHandler(
    messagesModel model.MessagesModel,
    messageSendsModel model.MessageSendsModel,
    templatesModel model.MessageTemplatesModel,
    settingsModel model.NotificationSettingsModel,
) *TemplateMessageHandler {
    return &TemplateMessageHandler{
        logger:             zerolog.GetLogger(),
        messagesModel:      messagesModel,
        messageSendsModel:  messageSendsModel,
        templatesModel:     templatesModel,
        settingsModel:      settingsModel,
    }
}

func (h *TemplateMessageHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.MessageTemplateEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "template_id":   event.Template.TemplateID,
        "template_code": event.Template.TemplateCode,
        "user_id":      event.UserID,
        "action":       event.Action,
        "event_type":   event.Type,
    }
    h.logger.Info(ctx, "Processing template message event", fields)

    switch event.Type {
    case types.MessageTemplateCreated:
        return h.handleTemplateCreated(ctx, event)
    case types.MessageTemplateUpdated:
        return h.handleTemplateUpdated(ctx, event)
    case types.MessageTemplateDeleted:
        return h.handleTemplateDeleted(ctx, event)
    default:
        return h.handleTemplateEvent(ctx, event)
    }
}

func (h *TemplateMessageHandler) handleTemplateCreated(ctx context.Context, event types.MessageTemplateEvent) error {
    template := &model.MessageTemplates{
        Code:            event.Template.TemplateCode,
        Name:            "New Template",
        TitleTemplate:   "New Template",
        ContentTemplate: event.Template.Content,
        Type:           1, // System template
        Channels:       "1,2,3", // All channels
        Status:         1, // Enabled
        CreatedAt:      event.Timestamp,
        UpdatedAt:      event.Timestamp,
    }

    _, err := h.templatesModel.Insert(ctx, template)
    if err != nil {
        return err
    }

    // Create notification for template creation
    message := &model.Messages{
        UserId:      uint64(event.UserID),
        Title:       "Template Created",
        Content:     "New message template created: " + event.Template.TemplateCode,
        Type:        1, // System message
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
        TemplateId: sql.NullInt64{Int64: event.Template.TemplateID, Valid: true},
        UserId:     uint64(event.UserID),
        Channel:    1, // System notification
        Status:     1, // Pending
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}

func (h *TemplateMessageHandler) handleTemplateUpdated(ctx context.Context, event types.MessageTemplateEvent) error {
    template, err := h.templatesModel.FindOneByCode(ctx, event.Template.TemplateCode)
    if err != nil {
        return err
    }

    template.ContentTemplate = event.Template.Content
    template.UpdatedAt = event.Timestamp

    if err := h.templatesModel.Update(ctx, template); err != nil {
        return err
    }

    // Create notification for template update
    message := &model.Messages{
        UserId:      uint64(event.UserID),
        Title:       "Template Updated",
        Content:     "Message template updated: " + event.Template.TemplateCode,
        Type:        1, // System message
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

    send := &model.MessageSends{
        MessageId:  uint64(messageId),
        TemplateId: sql.NullInt64{Int64: event.Template.TemplateID, Valid: true},
        UserId:     uint64(event.UserID),
        Channel:    1, // System notification
        Status:     1, // Pending
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}

func (h *TemplateMessageHandler) handleTemplateDeleted(ctx context.Context, event types.MessageTemplateEvent) error {
    if err := h.templatesModel.Delete(ctx, uint64(event.Template.TemplateID)); err != nil {
        return err
    }

    // Create notification for template deletion
    message := &model.Messages{
        UserId:      uint64(event.UserID),
        Title:       "Template Deleted",
        Content:     "Message template deleted: " + event.Template.TemplateCode,
        Type:        1, // System message
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

    send := &model.MessageSends{
        MessageId:  uint64(messageId),
        TemplateId: sql.NullInt64{Int64: event.Template.TemplateID, Valid: true},
        UserId:     uint64(event.UserID),
        Channel:    1, // System notification
        Status:     1, // Pending
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}

func (h *TemplateMessageHandler) handleTemplateEvent(ctx context.Context, event types.MessageTemplateEvent) error {
    // Handle other template events
    message := &model.Messages{
        UserId:      uint64(event.UserID),
        Title:       "Template Event",
        Content:     "Template event occurred: " + string(event.Type),
        Type:        1, // System message
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

    send := &model.MessageSends{
        MessageId:  uint64(messageId),
        TemplateId: sql.NullInt64{Int64: event.Template.TemplateID, Valid: true},
        UserId:     uint64(event.UserID),
        Channel:    1, // System notification
        Status:     1, // Pending
        CreatedAt:  time.Now(),
        UpdatedAt:  time.Now(),
    }

    _, err = h.messageSendsModel.Insert(ctx, send)
    return err
}