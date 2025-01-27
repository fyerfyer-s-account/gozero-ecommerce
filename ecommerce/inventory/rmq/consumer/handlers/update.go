package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateHandler struct {
    logger            *zerolog.Logger
    stocksModel       model.StocksModel
    stockRecordsModel model.StockRecordsModel
}

func NewUpdateHandler(
    stocksModel model.StocksModel,
    stockRecordsModel model.StockRecordsModel,
) *UpdateHandler {
    return &UpdateHandler{
        logger:            zerolog.GetLogger(),
        stocksModel:       stocksModel,
        stockRecordsModel: stockRecordsModel,
    }
}

func (h *UpdateHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.InventoryEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "type":         event.Type,
        "warehouse_id": event.WarehouseID,
    }
    h.logger.Info(ctx, "Processing stock update event", fields)

    switch event.Type {
    case types.StockUpdated:
        return h.handleStockUpdated(ctx, msg.Body)
    case types.StockOutOfStock:
        return h.handleStockOutOfStock(ctx, msg.Body)
    case types.StockLowStock:
        return h.handleStockLowStock(ctx, msg.Body)
    default:
        return nil
    }
}

func (h *UpdateHandler) handleStockUpdated(ctx context.Context, data []byte) error {
    var event types.StockUpdatedEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    // Get current stock
    stock, err := h.stocksModel.FindOneBySkuIdWarehouseId(ctx, uint64(event.SkuID), uint64(event.WarehouseID))
    if err != nil {
        return err
    }

    return h.stocksModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        // Update stock quantity
        stock.Available = int64(event.NewQuantity)
        stock.Total = stock.Available + stock.Locked

        if err := h.stocksModel.Update(ctx, stock); err != nil {
            return err
        }

        // Create stock record
        recordType := int64(1) // In
        if event.NewQuantity < event.OldQuantity {
            recordType = 2 // Out
        }

        quantity := event.NewQuantity - event.OldQuantity
        if quantity < 0 {
            quantity = -quantity
        }

        _, err := h.stockRecordsModel.Insert(ctx, &model.StockRecords{
            SkuId:       uint64(event.SkuID),
            WarehouseId: uint64(event.WarehouseID),
            Type:        recordType,
            Quantity:    int64(quantity),
            Remark:      sql.NullString{String: event.Reason, Valid: true},
        })

        // Check if stock is below threshold after update
        if stock.Available <= stock.AlertQuantity {
            h.sendLowStockAlert(ctx, stock)
        }

        return err
    })
}

func (h *UpdateHandler) handleStockOutOfStock(ctx context.Context, data []byte) error {
    var event types.StockOutOfStockEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    // Log out of stock event
    h.logger.Warn(ctx, "Stock out of stock", map[string]interface{}{
        "sku_id":       event.SkuID,
        "warehouse_id": event.WarehouseID,
        "quantity":     event.Quantity,
        "reason":       event.Reason,
    })

    // Create stock record
    _, err := h.stockRecordsModel.Insert(ctx, &model.StockRecords{
        SkuId:       uint64(event.SkuID),
        WarehouseId: uint64(event.WarehouseID),
        Type:        5, // Out of stock alert
        Quantity:    int64(event.Quantity),
        Remark:      sql.NullString{String: fmt.Sprintf("Out of stock: %s", event.Reason), Valid: true},
    })

    return err
}

func (h *UpdateHandler) handleStockLowStock(ctx context.Context, data []byte) error {
    var event types.StockLowStockEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    // Log low stock event
    h.logger.Warn(ctx, "Stock below threshold", map[string]interface{}{
        "sku_id":       event.SkuID,
        "warehouse_id": event.WarehouseID,
        "quantity":     event.Quantity,
        "threshold":    event.Threshold,
    })

    // Create stock record
    _, err := h.stockRecordsModel.Insert(ctx, &model.StockRecords{
        SkuId:       uint64(event.SkuID),
        WarehouseId: uint64(event.WarehouseID),
        Type:        6, // Low stock alert
        Quantity:    int64(event.Quantity),
        Remark:      sql.NullString{String: fmt.Sprintf("Stock below threshold: %d/%d", event.Quantity, event.Threshold), Valid: true},
    })

    return err
}

func (h *UpdateHandler) sendLowStockAlert(ctx context.Context, stock *model.Stocks) {
    h.logger.Warn(ctx, "Low stock alert", map[string]interface{}{
        "sku_id":       stock.SkuId,
        "warehouse_id": stock.WarehouseId,
        "available":    stock.Available,
        "threshold":    stock.AlertQuantity,
    })
}