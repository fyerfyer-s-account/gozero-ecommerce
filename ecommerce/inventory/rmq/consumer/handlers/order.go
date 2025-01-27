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

type OrderHandler struct {
    logger          *zerolog.Logger
    stocksModel     model.StocksModel
    stockLocksModel model.StockLocksModel
    stockRecordsModel model.StockRecordsModel
}

func NewOrderHandler(
    stocksModel model.StocksModel,
    stockLocksModel model.StockLocksModel,
    stockRecordsModel model.StockRecordsModel,
) *OrderHandler {
    return &OrderHandler{
        logger:          zerolog.GetLogger(),
        stocksModel:     stocksModel,
        stockLocksModel: stockLocksModel,
        stockRecordsModel: stockRecordsModel,
    }
}

func (h *OrderHandler) Handle(ctx context.Context, msg amqp.Delivery) error {
    var event types.OrderEvent
    if err := json.Unmarshal(msg.Body, &event); err != nil {
        return err
    }

    fields := map[string]interface{}{
        "order_no": event.OrderNo,
        "type":     event.Type,
    }
    h.logger.Info(ctx, "Processing order event", fields)

    switch event.Type {
    case types.OrderCreated:
        return h.handleOrderCreated(ctx, msg.Body)
    case types.OrderCancelled:
        return h.handleOrderCancelled(ctx, msg.Body)
    case types.OrderPaid:
        return h.handleOrderPaid(ctx, msg.Body)
    case types.OrderRefunded:
        return h.handleOrderRefunded(ctx, msg.Body)
    default:
        return nil
    }
}

func (h *OrderHandler) handleOrderCreated(ctx context.Context, data []byte) error {
    var event types.OrderCreatedEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    // Lock stock for each item
    return h.stocksModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        for _, item := range event.Items {
            // Lock stock
            err := h.stocksModel.Lock(ctx, session, uint64(item.SkuID), 1, int64(item.Quantity)) // use default warehouse
            if err != nil {
                return err
            }

            // Create lock record
            _, err = h.stockLocksModel.Insert(ctx, &model.StockLocks{
                OrderNo:     event.OrderNo,
                SkuId:      uint64(item.SkuID),
                WarehouseId: 1, // Default warehouse
                Quantity:   int64(item.Quantity),
                Status:    1, // Locked
            })
            if err != nil {
                return err
            }

            // Create stock record
            _, err = h.stockRecordsModel.Insert(ctx, &model.StockRecords{
                SkuId:      uint64(item.SkuID),
                WarehouseId: 1,
                Type:       3, // Lock
                Quantity:   int64(item.Quantity),
                OrderNo:    sql.NullString{String: event.OrderNo, Valid: true},
            })
            if err != nil {
                return err
            }
        }
        return nil
    })
}

func (h *OrderHandler) handleOrderCancelled(ctx context.Context, data []byte) error {
    var event types.OrderCancelledEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    // Find locked stocks
    locks, err := h.stockLocksModel.FindByOrderNo(ctx, event.OrderNo)
    if err != nil {
        return err
    }

    // Unlock stocks
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
                SkuId:      lock.SkuId,
                WarehouseId: lock.WarehouseId,
                Type:       4, // Unlock
                Quantity:   lock.Quantity,
                OrderNo:    sql.NullString{String: event.OrderNo, Valid: true},
                Remark:     sql.NullString{String: event.Reason, Valid: true},
            })
            if err != nil {
                return err
            }
        }
        return nil
    })
}

func (h *OrderHandler) handleOrderPaid(ctx context.Context, data []byte) error {
    var event types.OrderPaidEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    // Find locked stocks
    locks, err := h.stockLocksModel.FindByOrderNo(ctx, event.OrderNo)
    if err != nil {
        return err
    }

    // Deduct stocks
    return h.stocksModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        for _, lock := range locks {
            // Deduct stock
            err := h.stocksModel.Deduct(ctx, session, lock.SkuId, lock.WarehouseId, lock.Quantity)
            if err != nil {
                return err
            }

            // Update lock status
            lock.Status = 3 // Deducted
            err = h.stockLocksModel.Update(ctx, lock)
            if err != nil {
                return err
            }

            // Create stock record
            _, err = h.stockRecordsModel.Insert(ctx, &model.StockRecords{
                SkuId:      lock.SkuId,
                WarehouseId: lock.WarehouseId,
                Type:       2, // Out
                Quantity:   lock.Quantity,
                OrderNo:    sql.NullString{String: event.OrderNo, Valid: true},
            })
            if err != nil {
                return err
            }
        }
        return nil
    })
}

func (h *OrderHandler) handleOrderRefunded(ctx context.Context, data []byte) error {
    var event types.OrderRefundedEvent
    if err := json.Unmarshal(data, &event); err != nil {
        return err
    }

    // Return stock
    return h.stocksModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
        for _, item := range event.Items {
            // Increment available stock
            err := h.stocksModel.IncrAvailable(ctx, uint64(item.SkuID), 1, int64(item.Quantity))
            if err != nil {
                return err
            }

            // Create stock record
            _, err = h.stockRecordsModel.Insert(ctx, &model.StockRecords{
                SkuId:       uint64(item.SkuID),
                WarehouseId: 1, // Default warehouse
                Type:        1, // In
                Quantity:    int64(item.Quantity),
                OrderNo:     sql.NullString{String: event.OrderNo, Valid: true},
                Remark:      sql.NullString{String: fmt.Sprintf("Order refund: %s", event.Reason), Valid: true},
            })
            if err != nil {
                return err
            }
        }
        return nil
    })
}