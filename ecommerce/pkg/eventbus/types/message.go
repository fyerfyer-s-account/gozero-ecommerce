package types

import "time"

type MessageEventType string

const (
    // System message events
    MessageEventSent     MessageEventType = "message.event.sent"
    MessageEventReceived MessageEventType = "message.event.received"
    MessageEventFailed   MessageEventType = "message.event.failed"
    
    // Template message events
    MessageTemplateCreated  MessageEventType = "message.template.created"
    MessageTemplateUpdated  MessageEventType = "message.template.updated"
    MessageTemplateDeleted  MessageEventType = "message.template.deleted"
    
    // Batch message events
    MessageBatchCreated    MessageEventType = "message.batch.created"
    MessageBatchProcessing MessageEventType = "message.batch.processing"
    MessageBatchCompleted  MessageEventType = "message.batch.completed"
    MessageBatchFailed     MessageEventType = "message.batch.failed"
    
    // Payment related message events
    MessagePaymentSuccess  MessageEventType = "message.payment.success"
    MessagePaymentFailed   MessageEventType = "message.payment.failed"
)

// MessageEvent represents the base message event structure
type MessageEvent struct {
    Type      MessageEventType `json:"type"`
    UserID    int64           `json:"user_id"`
    Timestamp time.Time       `json:"timestamp"`
}

// MessageTemplate represents a message template
type MessageTemplate struct {
    TemplateID   int64  `json:"template_id"`
    TemplateCode string `json:"template_code"`
    Content      string `json:"content"`
    Variables    map[string]string `json:"variables"`
}

// MessageEventSentEvent represents a sent message event
type MessageEventSentEvent struct {
    MessageEvent
    MessageID   string            `json:"message_id"`
    Channel     string            `json:"channel"`      // sms, email, push, etc
    Content     string            `json:"content"`
    Recipient   string            `json:"recipient"`
    Status      string            `json:"status"`
    Variables   map[string]string `json:"variables,omitempty"`
}

// MessageTemplateEvent represents template related events
type MessageTemplateEvent struct {
    MessageEvent
    Template MessageTemplate `json:"template"`
    Action   string         `json:"action"` // created, updated, deleted
}

// MessageBatchEvent represents batch message events
type MessageBatchEvent struct {
    MessageEvent
    BatchID     string      `json:"batch_id"`
    Total       int32       `json:"total"`
    Completed   int32       `json:"completed"`
    Failed      int32       `json:"failed"`
    Status      string      `json:"status"`
}

// MessagePaymentSuccessEvent represents payment success notification
type MessagePaymentSuccessEvent struct {
    MessageEvent
    OrderNo    string  `json:"order_no"`
    PaymentNo  string  `json:"payment_no"`
    Amount     float64 `json:"amount"`
    Channel    string  `json:"channel"`
    TemplateID int64   `json:"template_id"`
}

// MessagePaymentFailedEvent represents payment failure notification
type MessagePaymentFailedEvent struct {
    MessageEvent
    OrderNo    string  `json:"order_no"`
    PaymentNo  string  `json:"payment_no"`
    Amount     float64 `json:"amount"`
    Reason     string  `json:"reason"`
    Channel    string  `json:"channel"`
    TemplateID int64   `json:"template_id"`
}

// Add validation methods
func (e *MessageEventSentEvent) Validate() error {
    if e.MessageID == "" || e.Channel == "" || e.Recipient == "" {
        return &NonRetryableError{
            EventError: &EventError{
                Code:    "INVALID_MESSAGE_EVENT",
                Message: "message_id, channel and recipient are required",
            },
        }
    }
    return nil
}

func (e *MessageTemplateEvent) Validate() error {
    if e.Template.TemplateID == 0 || e.Template.TemplateCode == "" {
        return &NonRetryableError{
            EventError: &EventError{
                Code:    "INVALID_TEMPLATE_EVENT",
                Message: "template_id and template_code are required",
            },
        }
    }
    return nil
}

func (e *MessageBatchEvent) Validate() error {
    if e.BatchID == "" {
        return &NonRetryableError{
            EventError: &EventError{
                Code:    "INVALID_BATCH_EVENT",
                Message: "batch_id is required",
            },
        }
    }
    return nil
}

func (e *MessagePaymentSuccessEvent) Validate() error {
    if e.OrderNo == "" || e.PaymentNo == "" || e.Channel == "" {
        return &NonRetryableError{
            EventError: &EventError{
                Code:    "INVALID_PAYMENT_SUCCESS_MESSAGE",
                Message: "order_no, payment_no and channel are required",
            },
        }
    }
    return nil
}

func (e *MessagePaymentFailedEvent) Validate() error {
    if e.OrderNo == "" || e.PaymentNo == "" || e.Channel == "" {
        return &NonRetryableError{
            EventError: &EventError{
                Code:    "INVALID_PAYMENT_FAILED_MESSAGE",
                Message: "order_no, payment_no and channel are required",
            },
        }
    }
    return nil
}