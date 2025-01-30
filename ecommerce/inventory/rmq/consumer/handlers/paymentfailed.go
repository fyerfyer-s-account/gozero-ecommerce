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

type PaymentFailedHandler struct {
    logger          *zerolog.Logger
    stocksModel     model.StocksModel
    stockLocksModel model.StockLocksModel
    stockRecordsModel model.StockRecordsModel
}

func NewPaymentFailedHandler(
    stocksModel model.StocksModel,
    stockLocksModel model.StockLocksModel,
    stockRecordsModel model.StockRecordsModel,
) *PaymentFailedHandler {
    return &PaymentFailedHandler{
        logger:           zerolog.GetLogger(),
        stocksModel:      stocksModel,
        stockLocksModel:  stockLocksModel,
        stockRecordsModel: stockRecordsModel,
    }
}

func (h *PaymentFailedHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.InventoryPaymentFailedEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "order_no":   event.OrderNo,
        "payment_no": event.PaymentNo,
        "items":      event.Items,
    }
    h.logger.Info(ctx, "Processing payment failed event", fields)

    // Start transaction
    err := h.stocksModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        // Get stock locks for order
        locks, err := h.stockLocksModel.FindAndLockByOrderNo(ctx, session, event.OrderNo)
        if err != nil {
            h.logger.Error(ctx, "Failed to find stock locks", err, fields)
            return err
        }

        // Release each locked stock
        for _, lock := range locks {
            // Unlock stock
            if err := h.stocksModel.Unlock(ctx, session, lock.SkuId, lock.WarehouseId, lock.Quantity); err != nil {
                h.logger.Error(ctx, "Failed to unlock stock", err, fields)
                return err
            }

            // Create stock record
            record := &model.StockRecords{
                SkuId:       lock.SkuId,
                WarehouseId: lock.WarehouseId,
                Type:        4, // Stock unlock
                Quantity:    lock.Quantity,
                OrderNo:     sql.NullString{String: event.OrderNo, Valid: true},
                Remark:      sql.NullString{String: event.Reason, Valid: true},
                Operator:    sql.NullString{String: "system", Valid: true},
            }
            
            if _, err := h.stockRecordsModel.Insert(ctx, record); err != nil {
                h.logger.Error(ctx, "Failed to create stock record", err, fields)
                return err
            }
        }

        // Update stock lock status
        if err := h.stockLocksModel.UpdateStatus(ctx, event.OrderNo, 1, 2); err != nil { // 1: locked -> 2: unlocked
            h.logger.Error(ctx, "Failed to update stock lock status", err, fields)
            return err
        }

        return nil
    })

    if err != nil {
        return err
    }

    h.logger.Info(ctx, "Successfully processed payment failed event", fields)
    return nil
}