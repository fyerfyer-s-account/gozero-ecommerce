package types

import "time"

type OrderEventType string

const (
	OrderCreated   OrderEventType = "order.created"
	OrderPaid      OrderEventType = "order.paid"
	OrderCancelled OrderEventType = "order.cancelled"
	OrderShipped   OrderEventType = "order.shipped"
	OrderCompleted OrderEventType = "order.completed"
	OrderRefunded  OrderEventType = "order.refunded"
)

// OrderEvent represents the base order event structure
type OrderEvent struct {
	Type      OrderEventType `json:"type"`
	OrderNo   string         `json:"order_no"`
	UserID    int64          `json:"user_id"`
	Timestamp time.Time      `json:"timestamp"`
}

// OrderCreatedEvent represents order creation event
type OrderCreatedEvent struct {
	OrderEvent
	Items       []OrderItem `json:"items"`
	TotalAmount float64     `json:"total_amount"`
	PayAmount   float64     `json:"pay_amount"`
	Address     string      `json:"address"`
	Receiver    string      `json:"receiver"`
	Phone       string      `json:"phone"`
}

// OrderItem represents an item in the order
type OrderItem struct {
	ProductID int64   `json:"product_id"`
	SkuID     int64   `json:"sku_id"`
	Quantity  int32   `json:"quantity"`
	Price     float64 `json:"price"`
}

// OrderPaidEvent represents order payment event
type OrderPaidEvent struct {
	OrderEvent
	PaymentNo     string  `json:"payment_no"`
	PaymentMethod int32   `json:"payment_method"`
	PayAmount     float64 `json:"pay_amount"`
}

// OrderCancelledEvent represents order cancellation event
type OrderCancelledEvent struct {
	OrderEvent
	Reason string `json:"reason"`
}

// OrderShippedEvent represents order shipping event
type OrderShippedEvent struct {
	OrderEvent
	ShippingNo string `json:"shipping_no"`
	Company    string `json:"company"`
}

// OrderCompletedEvent represents order completion event
type OrderCompletedEvent struct {
	OrderEvent
	ReceiveTime time.Time `json:"receive_time"`
}

type OrderAlertEvent struct {
	OrderEvent
	AlertType  string `json:"alert_type"`
	AlertLevel string `json:"alert_level"`
	Message    string `json:"message"`
}

// OrderStatusEventType represents the specific type of status change
type OrderStatusEventType string

const (
	OrderStatusPaid      OrderStatusEventType = "order.status.paid"
	OrderStatusShipped   OrderStatusEventType = "order.status.shipped"
	OrderStatusReceived  OrderStatusEventType = "order.status.received"
	OrderStatusCanceled  OrderStatusEventType = "order.status.canceled"
	OrderStatusRefunding OrderStatusEventType = "order.status.refunding"
)

// Enhanced OrderStatusChangedEvent
type OrderStatusChangedEvent struct {
	OrderEvent
	OldStatus  int32                `json:"old_status"`
	NewStatus  int32                `json:"new_status"`
	EventType  OrderStatusEventType `json:"event_type"`
	PaymentNo  string               `json:"payment_no,omitempty"`
	ShippingNo string               `json:"shipping_no,omitempty"`
	RefundNo   string               `json:"refund_no,omitempty"`
	Reason     string               `json:"reason,omitempty"`
}

// OrderRefundedEvent represents order refund event
type OrderRefundedEvent struct {
    OrderEvent
    Items      []OrderItem `json:"items"`
    RefundNo   string      `json:"refund_no"`
    RefundType int32       `json:"refund_type"`   
    Amount     float64     `json:"amount"`
    Reason     string      `json:"reason"`
}

// OrderPaymentSuccessEvent represents successful payment notification for order
type OrderPaymentSuccessEvent struct {
    OrderEvent
    PaymentNo     string    `json:"payment_no"`
    PaymentMethod int32     `json:"payment_method"`
    Amount        float64   `json:"amount"`
    PaidTime      time.Time `json:"paid_time"`
}

// OrderPaymentFailedEvent represents failed payment notification for order
type OrderPaymentFailedEvent struct {
    OrderEvent
    PaymentNo  string  `json:"payment_no"`
    Amount     float64 `json:"amount"`
    Reason     string  `json:"reason"`
    ErrorCode  string  `json:"error_code"`
}

// OrderPaymentRefundedEvent represents refund notification for order
type OrderPaymentRefundedEvent struct {
    OrderEvent
    PaymentNo    string    `json:"payment_no"`
    RefundNo     string    `json:"refund_no"`
    RefundAmount float64   `json:"refund_amount"`
    Reason       string    `json:"reason"`
    RefundTime   time.Time `json:"refund_time"`
}