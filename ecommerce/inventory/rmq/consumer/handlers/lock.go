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

type LockHandler struct {
    logger            *zerolog.Logger
    stocksModel       model.StocksModel
    stockLocksModel   model.StockLocksModel
    stockRecordsModel model.StockRecordsModel
}

func NewLockHandler(
    stocksModel model.StocksModel,
    stockLocksModel model.StockLocksModel,
    stockRecordsModel model.StockRecordsModel,
) *LockHandler {
    return &LockHandler{
        logger:            zerolog.GetLogger(),
        stocksModel:       stocksModel,
        stockLocksModel:   stockLocksModel,
        stockRecordsModel: stockRecordsModel,
    }
}

func (h *LockHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.InventoryEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "type":         event.Type,
        "warehouse_id": event.WarehouseID,
    }
    h.logger.Info(ctx, "Processing stock lock event", fields)

    switch event.Type {
    case types.StockLocked:
        return h.handleStockLocked(ctx, msg.Body)
    case types.StockUnlocked:
        return h.handleStockUnlocked(ctx, msg.Body)
    default:
        return nil
    }
}

func (h *LockHandler) handleStockLocked(ctx context.Context, data []byte) error {
    var event types.StockLockedEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    return h.stocksModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        for _, item := range event.Items {
            // Lock stock
            err := h.stocksModel.Lock(ctx, session, uint64(item.SkuID), uint64(event.WarehouseID), int64(item.Quantity))
            if err != nil {
                return err
            }

            // Create lock record
            _, err = h.stockLocksModel.Insert(ctx, &model.StockLocks{
                OrderNo:     event.OrderNo,
                SkuId:       uint64(item.SkuID),
                WarehouseId: uint64(event.WarehouseID),
                Quantity:    int64(item.Quantity),
                Status:      1, // Locked
            })
            if err != nil {
                return err
            }

            // Create stock record
            _, err = h.stockRecordsModel.Insert(ctx, &model.StockRecords{
                SkuId:       uint64(item.SkuID),
                WarehouseId: uint64(event.WarehouseID),
                Type:        3, // Lock
                Quantity:    int64(item.Quantity),
                OrderNo:     sql.NullString{String: event.OrderNo, Valid: true},
            })
            if err != nil {
                return err
            }
        }
        return nil
    })
}

func (h *LockHandler) handleStockUnlocked(ctx context.Context, data []byte) error {
    var event types.StockUnlockedEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    locks, err := h.stockLocksModel.FindByOrderNo(ctx, event.OrderNo)
    if err != nil {
        return err
    }

    return h.stocksModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        for _, lock := range locks {
            // Unlock stock
            err := h.stocksModel.Unlock(ctx, session, lock.SkuId, lock.WarehouseId, lock.Quantity)
            if err != nil {
                return err
            }

            // Update lock status
            lock.Status = 2 // Unlocked
            err = h.stockLocksModel.Update(ctx, lock)
            if err != nil {
                return err
            }

            // Create stock record
            _, err = h.stockRecordsModel.Insert(ctx, &model.StockRecords{
                SkuId:       lock.SkuId,
                WarehouseId: lock.WarehouseId,
                Type:        4, // Unlock
                Quantity:    lock.Quantity,
                OrderNo:     sql.NullString{String: event.OrderNo, Valid: true},
            })
            if err != nil {
                return err
            }
        }
        return nil
    })
}