package handlers

import (
    "context"
    "encoding/json"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type OrderHandler struct {
    logger         *zerolog.Logger
    cartItemsModel model.CartItemsModel
    cartStatsModel model.CartStatisticsModel
}

func NewOrderHandler(
    cartItemsModel model.CartItemsModel,
    cartStatsModel model.CartStatisticsModel,
) *OrderHandler {
    return &OrderHandler{
        logger:         zerolog.GetLogger(),
        cartItemsModel: cartItemsModel,
        cartStatsModel: cartStatsModel,
    }
}

func (h *OrderHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.OrderEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return &types.NonRetryableError{
            EventError: &types.EventError{
                Code:    "INVALID_EVENT_FORMAT",
                Message: "Failed to unmarshal order event",
                Err:     err,
            },
        }
    }

    fields := map[string]interface{}{
        "order_no": event.OrderNo,
        "user_id":  event.UserID,
        "type":     event.Type,
    }
    h.logger.Info(ctx, "Processing order event", fields)

    switch event.Type {
    case types.OrderCreated:
        return h.handleOrderCreated(ctx, msg.Body)
    default:
        return &types.NonRetryableError{
            EventError: &types.EventError{
                Code:    "UNKNOWN_EVENT_TYPE",
                Message: "Unknown order event type",
            },
        }
    }
}

func (h *OrderHandler) handleOrderCreated(ctx context.Context, data []byte) error {
    var event types.OrderCreatedEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "order_no": event.OrderNo,
        "user_id":  event.UserID,
        "items":    len(event.Items),
    }
    h.logger.Info(ctx, "Processing order created event", fields)

    // Get selected cart items for verification
    selectedItems, err := h.cartItemsModel.FindSelectedByUserId(ctx, uint64(event.UserID))
    if err != nil {
        h.logger.WithError(ctx, err, "Failed to get selected cart items", fields)
        return err
    }

    // Build map of ordered items
    orderedItems := make(map[int64]types.OrderItem)
    for _, item := range event.Items {
        orderedItems[item.SkuID] = item
    }

    // Delete ordered items from cart
    for _, cartItem := range selectedItems {
        if _, exists := orderedItems[int64(cartItem.SkuId)]; exists {
            if err := h.cartItemsModel.Delete(ctx, cartItem.Id); err != nil {
                h.logger.WithError(ctx, err, "Failed to delete cart item", fields)
                return err
            }
        }
    }

    // Recalculate cart statistics
    if err := h.cartStatsModel.RecalculateStats(ctx, uint64(event.UserID)); err != nil {
        h.logger.WithError(ctx, err, "Failed to recalculate cart statistics", fields)
        return err
    }

    h.logger.Info(ctx, "Successfully processed order created event", fields)
    return nil
}