package types

import (
	"encoding/json"
	"strconv"
	"time"
)

type EventType string

const (
	EventTypeOrderCreated   EventType = "order.created"
	EventTypeOrderPaid      EventType = "order.paid"
	EventTypeOrderCancelled EventType = "order.cancelled"
	EventTypeOrderShipped   EventType = "order.shipped"
	EventTypeOrderCompleted EventType = "order.completed"
	EventTypeRefundCreated  EventType = "order.refund.created"
	EventTypeRefundProcessed EventType = "order.refund.processed"
	EventTypeRefundCompleted EventType = "order.refund.completed"
)

type OrderEvent struct {
	ID        string      `json:"id"`
	Type      EventType   `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
	Metadata  Metadata    `json:"metadata"`
}

type Metadata struct {
	Source  string            `json:"source"`
	UserID  int64             `json:"userId,omitempty"`
	TraceID string            `json:"traceId"`
	Tags    map[string]string `json:"tags,omitempty"`
}

type OrderCreatedData struct {
	OrderNo     string      `json:"orderNo"`
	UserID      int64       `json:"userId"`
	TotalAmount float64     `json:"totalAmount"`
	Items       []OrderItem `json:"items"`
}

type OrderItem struct {
	SkuID     int64   `json:"skuId"`
	ProductID int64   `json:"productId"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
}

type OrderPaidData struct {
	OrderNo       string    `json:"orderNo"`
	PaymentNo     string    `json:"paymentNo"`
	PayAmount     float64   `json:"payAmount"`
	PaymentMethod int32     `json:"paymentMethod"`
	PayTime       time.Time `json:"payTime"` // Added to match OrderPaymentsModel
}

type OrderCancelledData struct {
    OrderNo string  `json:"orderNo"`
    OrderId int64   `json:"orderId"`
    Amount  float64 `json:"amount"`
    Reason  string  `json:"reason"`
}

type OrderShippedData struct {
    OrderNo    string `json:"orderNo"`
    OrderId    int64  `json:"orderId"`
    ShippingNo string `json:"shippingNo"`
    Company    string `json:"company"`
}

type RefundCreatedData struct {
    OrderNo     string  `json:"orderNo"`
    OrderId     int64   `json:"orderId"`
    RefundNo    string  `json:"refundNo"`
    Amount      float64 `json:"amount"`
    Reason      string  `json:"reason"`
    Images      string  `json:"images,omitempty"`
    Description string  `json:"description,omitempty"`
}

type RefundProcessedData struct {
    OrderNo     string    `json:"orderNo"`
    OrderId     int64     `json:"orderId"`
    RefundNo    string    `json:"refundNo"`
    Amount      float64   `json:"amount"`
    Status      int64     `json:"status"` // 1: approved, 2: rejected
    Reply       string    `json:"reply"`
    ProcessTime time.Time `json:"processTime"`
}

// Validate validates the event
func (e *OrderEvent) Validate() error {
    if e.ID == "" {
        return ErrEmptyEventID
    }
    if e.Type == "" {
        return ErrEmptyEventType
    }
    if e.Timestamp.IsZero() {
        return ErrEmptyTimestamp
    }
    if e.Data == nil {
        return ErrEmptyEventData
    }
    if e.Metadata.TraceID == "" {
        return ErrEmptyTraceID
    }
    return nil
}

// IsRetryable determines if the event can be retried
func (e *OrderEvent) IsRetryable() bool {
    switch e.Type {
    case EventTypeOrderCreated, 
         EventTypeOrderPaid,
         EventTypeRefundCreated,
         EventTypeRefundProcessed:
        return true
    default:
        return false
    }
}

// GetRetryCount gets the retry count from metadata
func (e *OrderEvent) GetRetryCount() int {
    if count, ok := e.Metadata.Tags["retry_count"]; ok {
        if v, err := strconv.Atoi(count); err == nil {
            return v
        }
    }
    return 0
}

func (e *OrderEvent) Marshal() ([]byte, error) {
    return json.Marshal(e)
}

// Unmarshal deserializes the event from JSON
func (e *OrderEvent) Unmarshal(data []byte) error {
    return json.Unmarshal(data, e)
}
