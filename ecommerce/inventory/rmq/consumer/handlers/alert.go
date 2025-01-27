package handlers

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
)

type AlertHandler struct {
	logger            *zerolog.Logger
	stocksModel       model.StocksModel
	stockRecordsModel model.StockRecordsModel
}

func NewAlertHandler(
	stocksModel model.StocksModel,
	stockRecordsModel model.StockRecordsModel,
) *AlertHandler {
	return &AlertHandler{
		logger:            zerolog.GetLogger(),
		stocksModel:       stocksModel,
		stockRecordsModel: stockRecordsModel,
	}
}

func (h *AlertHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
	var event types.StockAlertEvent
	if err := json.Unmarshal(msg.Body, &event); err != nil {
		return err
	}

	fields := map[string]interface{}{
		"sku_id":       event.SkuID,
		"warehouse_id": event.WarehouseID,
		"current":      event.Current,
		"threshold":    event.Threshold,
	}
	h.logger.Info(ctx, "Processing stock alert event", fields)

	// Get current stock
	stock, err := h.stocksModel.FindOneBySkuIdWarehouseId(ctx, uint64(event.SkuID), uint64(event.WarehouseID))
	if err != nil {
		return err
	}

	// Check if stock is out
	if stock.Available <= 0 {
		return h.handleOutOfStock(ctx, stock)
	}

	// Check if stock is below threshold
	if stock.Available <= stock.AlertQuantity {
		return h.handleLowStock(ctx, stock)
	}

	return nil
}

func (h *AlertHandler) handleOutOfStock(ctx context.Context, stock *model.Stocks) error {
	// Log out of stock event
	h.logger.Warn(ctx, "Stock out of stock", map[string]interface{}{
		"sku_id":       stock.SkuId,
		"warehouse_id": stock.WarehouseId,
		"available":    stock.Available,
	})

	// Create stock record for out of stock
	_, err := h.stockRecordsModel.Insert(ctx, &model.StockRecords{
		SkuId:       stock.SkuId,
		WarehouseId: stock.WarehouseId,
		Type:        5, // Out of stock alert
		Quantity:    stock.Available,
		Remark:      sql.NullString{String: "Stock out of stock alert", Valid: true},
	})

	return err
}

func (h *AlertHandler) handleLowStock(ctx context.Context, stock *model.Stocks) error {
	// Log low stock event
	h.logger.Warn(ctx, "Stock below threshold", map[string]interface{}{
		"sku_id":       stock.SkuId,
		"warehouse_id": stock.WarehouseId,
		"available":    stock.Available,
		"threshold":    stock.AlertQuantity,
	})

	// Create stock record for low stock
	_, err := h.stockRecordsModel.Insert(ctx, &model.StockRecords{
		SkuId:       stock.SkuId,
		WarehouseId: stock.WarehouseId,
		Type:        6, // Low stock alert
		Quantity:    stock.Available,
		Remark:      sql.NullString{String: "Stock below threshold alert", Valid: true},
	})

	return err
}
