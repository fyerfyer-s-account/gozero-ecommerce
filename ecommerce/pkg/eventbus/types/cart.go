package types

import "time"

type CartEventType string

const (
	CartUpdated        CartEventType = "cart.updated"
	CartCleared        CartEventType = "cart.cleared"
	CartSelected       CartEventType = "cart.selected"
	CartUnselected     CartEventType = "cart.unselected"
	CartPaymentSuccess CartEventType = "cart.payment.success"
)

// CartEvent represents the base cart event structure
type CartEvent struct {
	Type      CartEventType `json:"type"`
	UserID    int64         `json:"user_id"`
	Timestamp time.Time     `json:"timestamp"`
}

// CartItem represents an item in the cart
type CartItem struct {
	ProductID int64   `json:"product_id"`
	SkuID     int64   `json:"sku_id"`
	Quantity  int32   `json:"quantity"`
	Selected  bool    `json:"selected"`
	Price     float64 `json:"price"`
}

// CartUpdatedEvent represents cart update event
type CartUpdatedEvent struct {
	CartEvent
	Items []CartItem `json:"items"`
}

// CartClearedEvent represents cart clearing event
type CartClearedEvent struct {
	CartEvent
	Reason string `json:"reason"`
}

// CartSelectionEvent represents item selection status change
type CartSelectionEvent struct {
	CartEvent
	Items []CartItem `json:"items"`
}

// CartPaymentSuccessEvent represents cart clearing after successful payment
type CartPaymentSuccessEvent struct {
    CartEvent
    OrderNo    string    `json:"order_no"`
    PaymentNo  string    `json:"payment_no"`
    ClearTime  time.Time `json:"clear_time"`
    Items      []CartItem `json:"items"`
}

// Add validation method
func (e *CartPaymentSuccessEvent) Validate() error {
    if e.OrderNo == "" || e.PaymentNo == "" {
        return &NonRetryableError{
            EventError: &EventError{
                Code:    "INVALID_CART_PAYMENT_SUCCESS",
                Message: "order_no and payment_no are required",
            },
        }
    }
    return nil
}