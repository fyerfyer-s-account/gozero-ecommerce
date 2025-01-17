package logic

import (
	"context"
	"database/sql"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
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
			// Deduct locked stock
			err := l.svcCtx.StocksModel.Deduct(ctx, session, lock.SkuId, lock.WarehouseId, lock.Quantity)
			if err != nil {
				return err
			}

			// Create stock record
			record := &model.StockRecords{
				SkuId:       lock.SkuId,
				WarehouseId: lock.WarehouseId,
				Type:        2, // Stock out
				Quantity:    lock.Quantity,
				OrderNo:     sql.NullString{String: in.OrderNo, Valid: true},
				Remark:      sql.NullString{String: "stock_deduct", Valid: true},
			}
			_, err = l.svcCtx.StockRecordsModel.Insert(ctx, record)
			if err != nil {
				return err
			}

			// Publish stock update event
			if l.svcCtx.Producer != nil {
				err = l.svcCtx.Producer.PublishStockUpdate(
					ctx,
					&types.StockUpdateData{
						SkuID:       lock.SkuId,
						WarehouseID: lock.WarehouseId,  // Add this line
						Quantity:    -int32(lock.Quantity),
						Remark:      "stock_deduct",
					},
					0,
				)
				if err != nil {
					l.Logger.Errorf("Failed to publish stock update message: %v", err)
				}
			}
		}

		// Delete lock records
		return l.svcCtx.StockLocksModel.DeleteByOrderNo(ctx, in.OrderNo)
	})

	if err != nil {
		return nil, err
	}

	return &inventory.DeductStockResponse{
		Success: true,
	}, nil
}
