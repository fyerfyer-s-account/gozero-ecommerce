package handlers

import (
    "context"
    "encoding/json"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type StatusHandler struct {
    logger         *zerolog.Logger
    cartItemsModel model.CartItemsModel
    cartStatsModel model.CartStatisticsModel
}

func NewStatusHandler(
    cartItemsModel model.CartItemsModel,
    cartStatsModel model.CartStatisticsModel,
) *StatusHandler {
    return &StatusHandler{
        logger:         zerolog.GetLogger(),
        cartItemsModel: cartItemsModel,
        cartStatsModel: cartStatsModel,
    }
}

func (h *StatusHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.CartEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return &types.NonRetryableError{
            EventError: &types.EventError{
                Code:    "INVALID_EVENT_FORMAT",
                Message: "Failed to unmarshal cart status event",
                Err:     err,
            },
        }
    }

    fields := map[string]interface{}{
        "user_id": event.UserID,
        "type":    event.Type,
    }
    h.logger.Info(ctx, "Processing cart status event", fields)

    switch event.Type {
    case types.CartUpdated:
        return h.handleCartUpdated(ctx, msg.Body)
    case types.CartCleared:
        return h.handleCartCleared(ctx, msg.Body)
    case types.CartSelected:
        return h.handleCartSelected(ctx, msg.Body)
    case types.CartUnselected:
        return h.handleCartUnselected(ctx, msg.Body)
    default:
        return &types.NonRetryableError{
            EventError: &types.EventError{
                Code:    "UNKNOWN_EVENT_TYPE",
                Message: "Unknown cart event type",
            },
        }
    }
}

func (h *StatusHandler) handleCartUpdated(ctx context.Context, data []byte) error {
    var event types.CartUpdatedEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "user_id": event.UserID,
        "items":   len(event.Items),
    }
    h.logger.Info(ctx, "Processing cart update event", fields)

    // Update or insert cart items
    for _, item := range event.Items {
        cartItem := &model.CartItems{
            UserId:      uint64(event.UserID),
            ProductId:   uint64(item.ProductID),
            SkuId:      uint64(item.SkuID),
            Quantity:   int64(item.Quantity),
            Selected:   0,
        }
        if item.Selected {
            cartItem.Selected = 1
        }

        if _, err := h.cartItemsModel.Insert(ctx, cartItem); err != nil {
            h.logger.WithError(ctx, err, "Failed to update cart item", fields)
            return err
        }
    }

    // Recalculate cart statistics
    return h.cartStatsModel.RecalculateStats(ctx, uint64(event.UserID))
}

func (h *StatusHandler) handleCartCleared(ctx context.Context, data []byte) error {
    var event types.CartClearedEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "user_id": event.UserID,
        "reason":  event.Reason,
    }
    h.logger.Info(ctx, "Processing cart clear event", fields)

    // Delete all cart items
    if err := h.cartItemsModel.DeleteByUserId(ctx, uint64(event.UserID)); err != nil {
        h.logger.WithError(ctx, err, "Failed to clear cart items", fields)
        return err
    }

    // Reset cart statistics
    stats := &model.CartStatistics{
        UserId:           uint64(event.UserID),
        TotalQuantity:    0,
        SelectedQuantity: 0,
        TotalAmount:      0,
        SelectedAmount:   0,
    }
    return h.cartStatsModel.Upsert(ctx, stats)
}

func (h *StatusHandler) handleCartSelected(ctx context.Context, data []byte) error {
    var event types.CartSelectionEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "user_id": event.UserID,
        "items":   len(event.Items),
    }
    h.logger.Info(ctx, "Processing cart selection event", fields)

    // Update selected status
    for _, item := range event.Items {
        err := h.cartItemsModel.UpdateSelected(ctx, 
            uint64(event.UserID), 
            uint64(item.ProductID),
            uint64(item.SkuID),
            1,
        )
        if err != nil {
            h.logger.WithError(ctx, err, "Failed to update cart item selection", fields)
            return err
        }
    }

    // Recalculate cart statistics
    return h.cartStatsModel.RecalculateStats(ctx, uint64(event.UserID))
}

func (h *StatusHandler) handleCartUnselected(ctx context.Context, data []byte) error {
    var event types.CartSelectionEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "user_id": event.UserID,
        "items":   len(event.Items),
    }
    h.logger.Info(ctx, "Processing cart unselection event", fields)

    // Update unselected status
    for _, item := range event.Items {
        err := h.cartItemsModel.UpdateSelected(ctx, 
            uint64(event.UserID), 
            uint64(item.ProductID),
            uint64(item.SkuID),
            0,
        )
        if err != nil {
            h.logger.WithError(ctx, err, "Failed to update cart item unselection", fields)
            return err
        }
    }

    // Recalculate cart statistics
    return h.cartStatsModel.RecalculateStats(ctx, uint64(event.UserID))
}