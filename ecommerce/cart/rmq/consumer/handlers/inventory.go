package handlers

import (
    "context"
    "encoding/json"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type InventoryHandler struct {
    logger         *zerolog.Logger
    cartItemsModel model.CartItemsModel
    cartStatsModel model.CartStatisticsModel
}

func NewInventoryHandler(
    cartItemsModel model.CartItemsModel,
    cartStatsModel model.CartStatisticsModel,
) *InventoryHandler {
    return &InventoryHandler{
        logger:         zerolog.GetLogger(),
        cartItemsModel: cartItemsModel,
        cartStatsModel: cartStatsModel,
    }
}

func (h *InventoryHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.InventoryEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return &types.NonRetryableError{
            EventError: &types.EventError{
                Code:    "INVALID_EVENT_FORMAT",
                Message: "Failed to unmarshal inventory event",
                Err:     err,
            },
        }
    }

    fields := map[string]interface{}{
        "event_type":    event.Type,
        "warehouse_id":  event.WarehouseID,
        "timestamp":     event.Timestamp,
    }
    h.logger.Info(ctx, "Processing inventory event", fields)

    switch event.Type {
    case types.StockOutOfStock:
        return h.handleStockOutOfStock(ctx, msg.Body)
    case types.StockLowStock:
        return h.handleStockLowStock(ctx, msg.Body)
    case types.StockUpdated:
        return h.handleStockUpdated(ctx, msg.Body)
    default:
        return &types.NonRetryableError{
            EventError: &types.EventError{
                Code:    "UNKNOWN_EVENT_TYPE",
                Message: "Unknown inventory event type",
            },
        }
    }
}

func (h *InventoryHandler) handleStockOutOfStock(ctx context.Context, data []byte) error {
    var event types.StockOutOfStockEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "sku_id":   event.SkuID,
        "quantity": event.Quantity,
        "reason":   event.Reason,
    }
    h.logger.Info(ctx, "Processing out of stock event", fields)

    // Get all cart items containing this SKU
    items, err := h.cartItemsModel.FindByUserId(ctx, uint64(event.WarehouseID))
    if err != nil {
        h.logger.WithError(ctx, err, "Failed to find cart items", fields)
        return err
    }

    // Mark items as unavailable or remove them
    for _, item := range items {
        if item.SkuId == uint64(event.SkuID) {
            if err := h.cartItemsModel.Delete(ctx, item.Id); err != nil {
                h.logger.WithError(ctx, err, "Failed to delete cart item", fields)
                return err
            }
        }
    }

    // Recalculate cart statistics
    if err := h.cartStatsModel.RecalculateStats(ctx, uint64(event.WarehouseID)); err != nil {
        h.logger.WithError(ctx, err, "Failed to recalculate cart statistics", fields)
        return err
    }

    return nil
}

func (h *InventoryHandler) handleStockLowStock(ctx context.Context, data []byte) error {
    var event types.StockLowStockEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "sku_id":    event.SkuID,
        "quantity":  event.Quantity,
        "threshold": event.Threshold,
    }
    h.logger.Info(ctx, "Processing low stock event", fields)

    // Get all cart items containing this SKU
    items, err := h.cartItemsModel.FindByUserId(ctx, uint64(event.WarehouseID))
    if err != nil {
        h.logger.WithError(ctx, err, "Failed to find cart items", fields)
        return err
    }

    // Adjust quantities if necessary
    for _, item := range items {
        if item.SkuId == uint64(event.SkuID) && item.Quantity > int64(event.Quantity) {
            if err := h.cartItemsModel.UpdateQuantity(ctx, uint64(event.WarehouseID), item.ProductId, item.SkuId, int64(event.Quantity)); err != nil {
                h.logger.WithError(ctx, err, "Failed to update cart item quantity", fields)
                return err
            }
        }
    }

    // Recalculate cart statistics
    if err := h.cartStatsModel.RecalculateStats(ctx, uint64(event.WarehouseID)); err != nil {
        h.logger.WithError(ctx, err, "Failed to recalculate cart statistics", fields)
        return err
    }

    return nil
}

func (h *InventoryHandler) handleStockUpdated(ctx context.Context, data []byte) error {
    var event types.StockUpdatedEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "sku_id":        event.SkuID,
        "old_quantity":  event.OldQuantity,
        "new_quantity": event.NewQuantity,
        "reason":       event.Reason,
    }
    h.logger.Info(ctx, "Processing stock update event", fields)

    // Get all cart items containing this SKU
    items, err := h.cartItemsModel.FindByUserId(ctx, uint64(event.WarehouseID))
    if err != nil {
        h.logger.WithError(ctx, err, "Failed to find cart items", fields)
        return err
    }

    // Adjust quantities if necessary
    for _, item := range items {
        if item.SkuId == uint64(event.SkuID) && item.Quantity > int64(event.NewQuantity) {
            if err := h.cartItemsModel.UpdateQuantity(ctx, uint64(event.WarehouseID), item.ProductId, item.SkuId, int64(event.NewQuantity)); err != nil {
                h.logger.WithError(ctx, err, "Failed to update cart item quantity", fields)
                return err
            }
        }
    }

    // Recalculate cart statistics
    if err := h.cartStatsModel.RecalculateStats(ctx, uint64(event.WarehouseID)); err != nil {
        h.logger.WithError(ctx, err, "Failed to recalculate cart statistics", fields)
        return err
    }

    return nil
}