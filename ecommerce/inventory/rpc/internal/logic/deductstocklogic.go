package logic

import (
	"context"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type DeductStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeductStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeductStockLogic {
	return &DeductStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeductStockLogic) DeductStock(in *inventory.DeductStockRequest) (*inventory.DeductStockResponse, error) {
	if len(in.OrderNo) == 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Find locked stock records
	locks, err := l.svcCtx.StockLocksModel.FindByOrderNo(l.ctx, in.OrderNo)
	if err != nil {
		return nil, err
	}
	if len(locks) == 0 {
		return nil, zeroerr.ErrLockNotFound
	}

	err = l.svcCtx.StocksModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		for _, lock := range locks {
			// Get current stock before deduction
			stock, err := l.svcCtx.StocksModel.FindOneBySkuIdWarehouseId(ctx, lock.SkuId, lock.WarehouseId)
			if err != nil {
				return err
			}
			oldQuantity := int32(stock.Available + stock.Locked)

			// ...existing deduct and record creation code...

			// After successful deduction, publish events
			if l.svcCtx.Producer != nil {
				// Stock deducted event
				deductEvent := &types.StockDeductedEvent{
					InventoryEvent: types.InventoryEvent{
						Type:        types.StockDeducted,
						WarehouseID: int64(lock.WarehouseId),
						Timestamp:   time.Now(),
					},
					OrderNo: in.OrderNo,
					Items: []types.StockItem{
						{
							SkuID:    int64(lock.SkuId),
							Quantity: int32(lock.Quantity),
						},
					},
				}
				if err := l.svcCtx.Producer.PublishStockDeducted(ctx, deductEvent); err != nil {
					l.Logger.Error("Failed to publish stock deduct event", err)
				}

				// Stock update event
				newQuantity := oldQuantity - int32(lock.Quantity)
				updateEvent := &types.StockUpdatedEvent{
					InventoryEvent: types.InventoryEvent{
						Type:        types.StockUpdated,
						WarehouseID: int64(lock.WarehouseId),
						Timestamp:   time.Now(),
					},
					SkuID:       int64(lock.SkuId),
					OldQuantity: oldQuantity,
					NewQuantity: newQuantity,
					Reason:      "stock_deduct",
				}
				if err := l.svcCtx.Producer.PublishStockUpdated(ctx, updateEvent); err != nil {
					l.Logger.Error("Failed to publish stock update event", err)
				}

				// Check if stock level is critical
				remainingStock := stock.Available - lock.Quantity
				if remainingStock <= 0 {
					outOfStockEvent := &types.StockOutOfStockEvent{
						InventoryEvent: types.InventoryEvent{
							Type:        types.StockOutOfStock,
							WarehouseID: int64(lock.WarehouseId),
							Timestamp:   time.Now(),
						},
						SkuID:    int64(lock.SkuId),
						Quantity: 0,
						Reason:   "Stock depleted after deduction",
					}
					if err := l.svcCtx.Producer.PublishStockOutOfStock(ctx, outOfStockEvent); err != nil {
						l.Logger.Error("Failed to publish stock out event", err)
					}
				} else if remainingStock <= stock.AlertQuantity {
					lowStockEvent := &types.StockLowStockEvent{
						InventoryEvent: types.InventoryEvent{
							Type:        types.StockLowStock,
							WarehouseID: int64(lock.WarehouseId),
							Timestamp:   time.Now(),
						},
						SkuID:     int64(lock.SkuId),
						Quantity:  int32(remainingStock),
						Threshold: int32(stock.AlertQuantity),
					}
					if err := l.svcCtx.Producer.PublishStockLowStock(ctx, lowStockEvent); err != nil {
						l.Logger.Error("Failed to publish low stock event", err)
					}
				}
			}
		}

		// ...existing delete locks code...
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &inventory.DeductStockResponse{
		Success: true,
	}, nil
}
