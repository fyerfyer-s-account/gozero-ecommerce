package types

import "time"

type PaymentEventType string

const (
    PaymentCreated  PaymentEventType = "payment.created"
    PaymentSuccess  PaymentEventType = "payment.success"
    PaymentFailed   PaymentEventType = "payment.failed"
    PaymentRefund   PaymentEventType = "payment.refund"
    PaymentVerified PaymentEventType = "payment.verified"
)

// PaymentEvent represents the base payment event structure
type PaymentEvent struct {
    Type      PaymentEventType `json:"type"`
    OrderNo   string           `json:"order_no"`
    PaymentNo string           `json:"payment_no"`
    Timestamp time.Time        `json:"timestamp"`
}

// PaymentCreatedEvent represents payment creation
type PaymentCreatedEvent struct {
    PaymentEvent
    Amount        float64 `json:"amount"`
    PaymentMethod int32   `json:"payment_method"`
    PayURL       string  `json:"pay_url"`
}

// PaymentSuccessEvent represents successful payment
type PaymentSuccessEvent struct {
    PaymentEvent
    Amount        float64   `json:"amount"`
    PaymentMethod int32     `json:"payment_method"`
    PaidTime     time.Time `json:"paid_time"`
}

// PaymentFailedEvent represents failed payment
type PaymentFailedEvent struct {
    PaymentEvent
    Amount  float64 `json:"amount"`
    Reason  string  `json:"reason"`
    ErrorCode string `json:"error_code"`
}

// PaymentRefundEvent represents refund event
type PaymentRefundEvent struct {
    PaymentEvent
    RefundNo     string    `json:"refund_no"`
    RefundAmount float64   `json:"refund_amount"`
    Reason      string    `json:"reason"`
    RefundTime  time.Time `json:"refund_time"`
}

// PaymentVerificationEvent represents payment verification
type PaymentVerificationEvent struct {
    PaymentEvent
    Verified bool   `json:"verified"`
    Message  string `json:"message"`
}