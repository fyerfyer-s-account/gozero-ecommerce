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

type UpdateStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateStockLogic {
	return &UpdateStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateStockLogic) UpdateStock(in *inventory.UpdateStockRequest) (*inventory.UpdateStockResponse, error) {
	// Input validation
	if in.SkuId <= 0 || in.WarehouseId <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	err := l.svcCtx.StocksModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// Get current stock
		stock, err := l.svcCtx.StocksModel.FindOneBySkuIdWarehouseId(ctx, uint64(in.SkuId), uint64(in.WarehouseId))
		if err != nil && err != model.ErrNotFound {
			return err
		}

		if err == model.ErrNotFound {
			if in.Quantity < 0 {
				return zeroerr.ErrInsufficientStock
			}
			// Create new stock
			stock = &model.Stocks{
				SkuId:       uint64(in.SkuId),
				WarehouseId: uint64(in.WarehouseId),
				Available:   int64(in.Quantity),
				Locked:      0,
				Total:       int64(in.Quantity),
			}
			_, err = l.svcCtx.StocksModel.Insert(ctx, stock)
		} else {
			// Update existing stock
			if in.Quantity > 0 {
				err = l.svcCtx.StocksModel.IncrAvailable(ctx, uint64(in.SkuId), uint64(in.WarehouseId), int64(in.Quantity))
			} else if in.Quantity < 0 {
				if stock.Available < -int64(in.Quantity) {
					return zeroerr.ErrInsufficientStock
				}
				err = l.svcCtx.StocksModel.DecrAvailable(ctx, uint64(in.SkuId), uint64(in.WarehouseId), -int64(in.Quantity))
			}
		}
		if err != nil {
			return err
		}

		// Create stock record
		record := &model.StockRecords{
			SkuId:       uint64(in.SkuId),
			WarehouseId: uint64(in.WarehouseId),
			Type:        1, // Stock update
			Quantity:    int64(in.Quantity),
			Remark:      sql.NullString{String: in.Remark, Valid: len(in.Remark) > 0},
		}
		_, err = l.svcCtx.StockRecordsModel.Insert(ctx, record)
		if err != nil {
			return err
		}

		 // Publish stock update event
		if l.svcCtx.Producer != nil {
			if err := l.svcCtx.Producer.PublishStockUpdate(
				ctx,
				&types.StockUpdateData{
					SkuID:       uint64(in.SkuId),
					WarehouseID: uint64(in.WarehouseId),
					Quantity:    in.Quantity,
					Remark:      record.Remark.String,
				},
				0,
			); err != nil {
				l.Logger.Errorf("Failed to publish stock update message: %v", err)
			}
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return &inventory.UpdateStockResponse{
		Success: true,
	}, nil
}
