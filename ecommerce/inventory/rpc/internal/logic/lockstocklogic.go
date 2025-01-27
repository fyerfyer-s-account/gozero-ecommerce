package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type LockStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLockStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LockStockLogic {
	return &LockStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 库存锁定/解锁
func (l *LockStockLogic) LockStock(in *inventory.LockStockRequest) (*inventory.LockStockResponse, error) {
	// Input validation
	if len(in.OrderNo) == 0 || len(in.Items) == 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Track failed items
	failedItems := make([]*inventory.LockFailedItem, 0)

	err := l.svcCtx.StocksModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// Process each item
		lockRecords := make([]*model.StockLocks, 0, len(in.Items))
		stockRecords := make([]*model.StockRecords, 0, len(in.Items))

		// Create rmq message in advance
		var rmqItems []types.StockItem

		for _, item := range in.Items {
			if item.SkuId <= 0 || item.Quantity <= 0 || item.WarehouseId <= 0 {
				failedItems = append(failedItems, &inventory.LockFailedItem{
					SkuId:  item.SkuId,
					Reason: "invalid parameters",
				})
				continue
			}

			// Try to lock stock
			err := l.svcCtx.StocksModel.Lock(ctx, session,
				uint64(item.SkuId),
				uint64(item.WarehouseId),
				int64(item.Quantity))

			if err != nil {
				failedItems = append(failedItems, &inventory.LockFailedItem{
					SkuId:  item.SkuId,
					Reason: "insufficient stock",
				})
				continue
			}

			// Create lock record
			lockRecord := &model.StockLocks{
				OrderNo:     in.OrderNo,
				SkuId:       uint64(item.SkuId),
				WarehouseId: uint64(item.WarehouseId),
				Quantity:    int64(item.Quantity),
				Status:      1, // Locked
			}
			lockRecords = append(lockRecords, lockRecord)

			// Create stock record
			stockRecord := &model.StockRecords{
				SkuId:       uint64(item.SkuId),
				WarehouseId: uint64(item.WarehouseId),
				Type:        3, // Lock
				Quantity:    int64(item.Quantity),
				OrderNo:     sql.NullString{String: in.OrderNo, Valid: true},
				Remark:      sql.NullString{String: "stock_lock", Valid: true},
			}
			stockRecords = append(stockRecords, stockRecord)
			rmqItems = append(rmqItems, types.StockItem{
				SkuID:    item.SkuId,
				Quantity: item.Quantity,
			})
		}

		// Batch insert lock records
		if len(lockRecords) > 0 {
			if err := l.svcCtx.StockLocksModel.BatchInsert(ctx, lockRecords); err != nil {
				return err
			}
		}

		// Batch insert stock records
		if len(stockRecords) > 0 {
			if err := l.svcCtx.StockRecordsModel.BatchInsert(ctx, stockRecords); err != nil {
				return err
			}
		}

		// Publish stock lock event if successful
		if l.svcCtx.Producer != nil {
			event := &types.StockLockedEvent{
				InventoryEvent: types.InventoryEvent{
					Type:        types.StockLocked,
					WarehouseID: int64(lockRecords[0].WarehouseId), // Use first record's warehouse
					Timestamp:   time.Now(),
				},
				OrderNo: in.OrderNo,
				Items:   rmqItems,
			}

			if err := l.svcCtx.Producer.PublishStockLocked(ctx, event); err != nil {
				l.Logger.Error("Failed to publish stock lock event", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &inventory.LockStockResponse{
		Success:     len(failedItems) == 0,
		FailedItems: failedItems,
	}, nil
}
