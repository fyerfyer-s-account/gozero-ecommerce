package types

import "time"

type CartEventType string

const (
	CartUpdated    CartEventType = "cart.updated"
	CartCleared    CartEventType = "cart.cleared"
	CartSelected   CartEventType = "cart.selected"
	CartUnselected CartEventType = "cart.unselected"
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
