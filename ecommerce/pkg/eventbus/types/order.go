package types

import "time"

type OrderEventType string

const (
	OrderCreated       OrderEventType = "order.created"
	OrderPaid          OrderEventType = "order.paid"
	OrderCancelled     OrderEventType = "order.cancelled"
	OrderShipped       OrderEventType = "order.shipped"
	OrderCompleted     OrderEventType = "order.completed"
	OrderRefunded      OrderEventType = "order.refunded"
	OrderStatusChanged OrderEventType = "order.status.changed"
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

// OrderStatusChangedEvent represents order status change event
type OrderStatusChangedEvent struct {
    OrderEvent
    OldStatus int32 `json:"old_status"`
    NewStatus int32 `json:"new_status"`
}