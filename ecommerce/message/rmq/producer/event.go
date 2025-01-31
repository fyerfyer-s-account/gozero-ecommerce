package producer

import (
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
)

// Event message events
func NewMessageEventSentEvent(
    messageId string,
    userId int64,
    channel string,
    content string,
    recipient string,
    status string,
    variables map[string]string,
) *types.MessageEventSentEvent {
    return &types.MessageEventSentEvent{
        MessageEvent: types.MessageEvent{
            Type:      types.MessageEventSent,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        MessageID: messageId,
        Channel:   channel,
        Content:   content,
        Recipient: recipient,
        Status:    status,
        Variables: variables,
    }
}

// Template message events
func NewMessageTemplateEvent(
    userId int64,
    template types.MessageTemplate,
    action string,
    eventType types.MessageEventType,
) *types.MessageTemplateEvent {
    return &types.MessageTemplateEvent{
        MessageEvent: types.MessageEvent{
            Type:      eventType,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        Template: template,
        Action:   action,
    }
}

// Batch message events
func NewMessageBatchEvent(
    userId int64,
    batchId string, 
    total int32,
    completed int32,
    failed int32,
    status string,
    eventType types.MessageEventType,
) *types.MessageBatchEvent {
    return &types.MessageBatchEvent{
        MessageEvent: types.MessageEvent{
            Type:      eventType,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        BatchID:   batchId,
        Total:     total,
        Completed: completed,
        Failed:    failed,
        Status:    status,
    }
}

// Payment success notification event
func NewMessagePaymentSuccessEvent(
    userId int64,
    orderNo string,
    paymentNo string,
    amount float64,
    channel string,
    templateId int64,
) *types.MessagePaymentSuccessEvent {
    return &types.MessagePaymentSuccessEvent{
        MessageEvent: types.MessageEvent{
            Type:      types.MessagePaymentSuccess,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        OrderNo:    orderNo,
        PaymentNo:  paymentNo,
        Amount:     amount,
        Channel:    channel,
        TemplateID: templateId,
    }
}

// Payment failed notification event
func NewMessagePaymentFailedEvent(
    userId int64,
    orderNo string,
    paymentNo string,
    amount float64,
    reason string,
    channel string,
    templateId int64,
) *types.MessagePaymentFailedEvent {
    return &types.MessagePaymentFailedEvent{
        MessageEvent: types.MessageEvent{
            Type:      types.MessagePaymentFailed,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        OrderNo:    orderNo,
        PaymentNo:  paymentNo,
        Amount:     amount,
        Reason:     reason,
        Channel:    channel,
        TemplateID: templateId,
    }
}