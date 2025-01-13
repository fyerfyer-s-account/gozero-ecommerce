package types

import (
	"time"
)

type EventType string

const (
	EventTypeOrderCreated   EventType = "order.created"
	EventTypeOrderPaid      EventType = "order.paid"
	EventTypeOrderCancelled EventType = "order.cancelled"
	EventTypeOrderShipped   EventType = "order.shipped"
	EventTypeOrderCompleted EventType = "order.completed"
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
