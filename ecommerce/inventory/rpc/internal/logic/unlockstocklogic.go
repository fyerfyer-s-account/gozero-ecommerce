package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UnlockStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUnlockStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnlockStockLogic {
	return &UnlockStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UnlockStockLogic) UnlockStock(in *inventory.UnlockStockRequest) (*inventory.UnlockStockResponse, error) {
    // Validate input
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

    err = l.svcCtx.StockLocksModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
        stockRecords := make([]*model.StockRecords, 0, len(locks))

        for _, lock := range locks {
            // Unlock stock
            err := l.svcCtx.StocksModel.Unlock(ctx, session, lock.SkuId, lock.WarehouseId, lock.Quantity)
            if err != nil {
                return zeroerr.ErrStockUnlockFailed
            }

            // Create stock record
            stockRecord := &model.StockRecords{
                SkuId:       lock.SkuId,
                WarehouseId: lock.WarehouseId,
                Type:        4, // Unlock
                Quantity:    lock.Quantity,
                OrderNo:     sql.NullString{String: in.OrderNo, Valid: true},
                Remark:      sql.NullString{String: "stock_unlock", Valid: true},
            }
            stockRecords = append(stockRecords, stockRecord)
        }

        // Batch insert stock records
        if len(stockRecords) > 0 {
            if err := l.svcCtx.StockRecordsModel.BatchInsert(ctx, stockRecords); err != nil {
                return err
            }
        }

        // Create and publish stock unlock event
        event := &types.StockUnlockedEvent{
            InventoryEvent: types.InventoryEvent{
                Type:        types.StockUnlocked,
                WarehouseID: int64(locks[0].WarehouseId), // Use first lock's warehouse
                Timestamp:   time.Now(),
            },
            OrderNo: in.OrderNo,
        }

        if err := l.svcCtx.Producer.PublishStockUnlocked(ctx, event); err != nil {
            l.Logger.Error("Failed to publish stock unlock event", err)
        }

        // Delete lock records
        return l.svcCtx.StockLocksModel.DeleteByOrderNo(ctx, in.OrderNo)
    })

    if err != nil {
        return nil, err
    }

    return &inventory.UnlockStockResponse{
        Success: true,
    }, nil
}