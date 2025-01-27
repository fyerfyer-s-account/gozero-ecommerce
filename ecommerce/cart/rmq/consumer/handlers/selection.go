package handlers

import (
    "context"
    "encoding/json"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type SelectionHandler struct {
    logger         *zerolog.Logger
    cartItemsModel model.CartItemsModel
    cartStatsModel model.CartStatisticsModel
}

func NewSelectionHandler(
    cartItemsModel model.CartItemsModel,
    cartStatsModel model.CartStatisticsModel,
) *SelectionHandler {
    return &SelectionHandler{
        logger:         zerolog.GetLogger(),
        cartItemsModel: cartItemsModel,
        cartStatsModel: cartStatsModel,
    }
}

func (h *SelectionHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.CartSelectionEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return &types.NonRetryableError{
            EventError: &types.EventError{
                Code:    "INVALID_EVENT_FORMAT",
                Message: "Failed to unmarshal cart selection event",
                Err:     err,
            },
        }
    }

    fields := map[string]interface{}{
        "user_id": event.UserID,
        "items":   len(event.Items),
    }
    h.logger.Info(ctx, "Processing cart selection event", fields)

    // Handle individual item selections
    for _, item := range event.Items {
        selected := int64(0)
        if item.Selected {
            selected = 1
        }

        if err := h.cartItemsModel.UpdateSelected(ctx, uint64(event.UserID), uint64(item.ProductID), uint64(item.SkuID), selected); err != nil {
            h.logger.WithError(ctx, err, "Failed to update cart item selection", map[string]interface{}{
                "user_id":    event.UserID,
                "product_id": item.ProductID,
                "sku_id":     item.SkuID,
                "selected":   selected,
            })
            return err
        }
    }

    // Recalculate cart statistics after selection changes
    if err := h.cartStatsModel.RecalculateStats(ctx, uint64(event.UserID)); err != nil {
        h.logger.WithError(ctx, err, "Failed to recalculate cart statistics", fields)
        return err
    }

    h.logger.Info(ctx, "Successfully processed cart selection event", fields)
    return nil
}

func (h *SelectionHandler) HandleAllSelection(ctx context.Context, msg amqp.Delivery) error {
    var event types.CartEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return &types.NonRetryableError{
            EventError: &types.EventError{
                Code:    "INVALID_EVENT_FORMAT",
                Message: "Failed to unmarshal cart selection event",
                Err:     err,
            },
        }
    }

    fields := map[string]interface{}{
        "user_id": event.UserID,
        "type":    event.Type,
    }
    h.logger.Info(ctx, "Processing cart all selection event", fields)

    // Determine selection status based on event type
    selected := int64(1)
    if event.Type == types.CartUnselected {
        selected = 0
    }

    // Update all items selection status
    if err := h.cartItemsModel.UpdateAllSelected(ctx, uint64(event.UserID), selected); err != nil {
        h.logger.WithError(ctx, err, "Failed to update all cart items selection", fields)
        return err
    }

    // Recalculate cart statistics
    if err := h.cartStatsModel.RecalculateStats(ctx, uint64(event.UserID)); err != nil {
        h.logger.WithError(ctx, err, "Failed to recalculate cart statistics", fields)
        return err
    }

    h.logger.Info(ctx, "Successfully processed cart all selection event", fields)
    return nil
}