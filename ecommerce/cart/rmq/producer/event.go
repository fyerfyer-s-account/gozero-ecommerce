package producer

import (
    "time"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
)

func NewCartUpdatedEvent(userId int64, items []types.CartItem) *types.CartUpdatedEvent {
    return &types.CartUpdatedEvent{
        CartEvent: types.CartEvent{
            Type:      types.CartUpdated,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        Items: items,
    }
}

func NewCartClearedEvent(userId int64, reason string) *types.CartClearedEvent {
    return &types.CartClearedEvent{
        CartEvent: types.CartEvent{
            Type:      types.CartCleared,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        Reason: reason,
    }
}

func NewCartSelectionEvent(userId int64, items []types.CartItem, selected bool) *types.CartSelectionEvent {
    eventType := types.CartSelected
    if !selected {
        eventType = types.CartUnselected
    }
    
    return &types.CartSelectionEvent{
        CartEvent: types.CartEvent{
            Type:      eventType,
            UserID:    userId,
            Timestamp: time.Now(),
        },
        Items: items,
    }
}