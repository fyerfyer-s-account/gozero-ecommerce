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

type CreateStockInLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateStockInLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateStockInLogic {
	return &CreateStockInLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 入库/出库
func (l *CreateStockInLogic) CreateStockIn(in *inventory.CreateStockInRequest) (*inventory.CreateStockInResponse, error) {
	// Input validation
	if in.WarehouseId <= 0 || len(in.Items) == 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Begin transaction
	err := l.svcCtx.StocksModel.Trans(l.ctx, func(ctx context.Context, session sqlx.Session) error {
		// Process each item
		for _, item := range in.Items {
			if item.SkuId <= 0 || item.Quantity <= 0 {
				return zeroerr.ErrInvalidParam
			}

			// Update stock
			stock, err := l.svcCtx.StocksModel.FindOneBySkuIdWarehouseId(l.ctx, uint64(item.SkuId), uint64(in.WarehouseId))
			if err != nil && err != model.ErrNotFound {
				return err
			}

			if err == model.ErrNotFound {
				// Create new stock
				stock = &model.Stocks{
					SkuId:       uint64(item.SkuId),
					WarehouseId: uint64(in.WarehouseId),
					Available:   int64(item.Quantity),
					Locked:      0,
					Total:       int64(item.Quantity),
				}
				_, err = l.svcCtx.StocksModel.Insert(ctx, stock)
				if err != nil {
					return err
				}
			} else {
				// Update existing stock
				err = l.svcCtx.StocksModel.IncrAvailable(ctx, uint64(item.SkuId), uint64(in.WarehouseId), int64(item.Quantity))
				if err != nil {
					return err
				}
			}

			// Create stock record
			record := &model.StockRecords{
				SkuId:       uint64(item.SkuId),
				WarehouseId: uint64(in.WarehouseId),
				Type:        1, // Stock in
				Quantity:    int64(item.Quantity),
				Remark:      sql.NullString{String: in.Remark, Valid: len(in.Remark) > 0},
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
						SkuID:       uint64(item.SkuId),
						WarehouseID: uint64(in.WarehouseId),
						Quantity:    item.Quantity,
						Remark:      "stock_in",
					},
					0, // Pass userId if available from context
				)
				if err != nil {
					l.Logger.Errorf("Failed to publish stock update message: %v", err)
				}
			}
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &inventory.CreateStockInResponse{
		Success: true,
	}, nil
}
