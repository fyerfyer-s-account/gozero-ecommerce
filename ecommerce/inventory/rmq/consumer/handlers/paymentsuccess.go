package handlers

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
	"github.com/streadway/amqp"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type PaymentSuccessHandler struct {
    logger           *zerolog.Logger
    stocksModel      model.StocksModel
    stockLocksModel  model.StockLocksModel
    stockRecordsModel model.StockRecordsModel
}

func NewPaymentSuccessHandler(
    stocksModel model.StocksModel,
    stockLocksModel model.StockLocksModel,
    stockRecordsModel model.StockRecordsModel,
) *PaymentSuccessHandler {
    return &PaymentSuccessHandler{
        logger:            zerolog.GetLogger(),
        stocksModel:       stocksModel,
        stockLocksModel:   stockLocksModel,
        stockRecordsModel: stockRecordsModel,
    }
}

func (h *PaymentSuccessHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.InventoryPaymentSuccessEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "order_no":   event.OrderNo,
        "payment_no": event.PaymentNo,
        "items":      event.Items,
    }
    h.logger.Info(ctx, "Processing payment success event", fields)

    // Start transaction
    err := h.stocksModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        // Get stock locks for order
        locks, err := h.stockLocksModel.FindAndLockByOrderNo(ctx, session, event.OrderNo)
        if err != nil {
            h.logger.Error(ctx, "Failed to find stock locks", err, fields)
            return err
        }

        // Process each locked stock
        for _, lock := range locks {
            // Deduct locked stock
            if err := h.stocksModel.Deduct(ctx, session, lock.SkuId, lock.WarehouseId, lock.Quantity); err != nil {
                h.logger.Error(ctx, "Failed to deduct stock", err, fields)
                return err
            }

            // Create stock record
            record := &model.StockRecords{
                SkuId:       lock.SkuId,
                WarehouseId: lock.WarehouseId,
                Type:        2, // Stock deduction
                Quantity:    lock.Quantity,
                OrderNo:     sql.NullString{String: event.OrderNo, Valid: true},
                Remark:      sql.NullString{String: "Payment success stock deduction", Valid: true},
                Operator:    sql.NullString{String: "system", Valid: true},
            }
            
            if _, err := h.stockRecordsModel.Insert(ctx, record); err != nil {
                h.logger.Error(ctx, "Failed to create stock record", err, fields)
                return err
            }
        }

        // Update stock lock status
        if err := h.stockLocksModel.UpdateStatus(ctx, event.OrderNo, 1, 3); err != nil { // 1: locked -> 3: deducted
            h.logger.Error(ctx, "Failed to update stock lock status", err, fields)
            return err
        }

        return nil
    })

    if err != nil {
        return err
    }

    h.logger.Info(ctx, "Successfully processed payment success event", fields)
    return nil
}