package handlers

import (
    "context"
    "encoding/json"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type ClearHandler struct {
    logger            *zerolog.Logger
    cartItemsModel    model.CartItemsModel
    cartStatsModel    model.CartStatisticsModel
}

func NewClearHandler(
    cartItemsModel model.CartItemsModel,
    cartStatsModel model.CartStatisticsModel,
) *ClearHandler {
    return &ClearHandler{
        logger:            zerolog.GetLogger(),
        cartItemsModel:    cartItemsModel,
        cartStatsModel:    cartStatsModel,
    }
}

func (h *ClearHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.CartClearedEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return &types.NonRetryableError{
            EventError: &types.EventError{
                Code:    "INVALID_EVENT_FORMAT",
                Message: "Failed to unmarshal cart clear event",
                Err:     err,
            },
        }
    }

    fields := map[string]interface{}{
        "user_id": event.UserID,
        "reason":  event.Reason,
    }
    h.logger.Info(ctx, "Processing cart clear event", fields)

    // Delete all cart items for the user
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

    if err := h.cartStatsModel.Upsert(ctx, stats); err != nil {
        h.logger.WithError(ctx, err, "Failed to reset cart statistics", fields)
        return err
    }

    h.logger.Info(ctx, "Successfully cleared cart", fields)
    return nil
}