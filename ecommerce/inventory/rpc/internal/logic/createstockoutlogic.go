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

type CreateStockOutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateStockOutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateStockOutLogic {
	return &CreateStockOutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateStockOutLogic) CreateStockOut(in *inventory.CreateStockOutRequest) (*inventory.CreateStockOutResponse, error) {
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

			// First check if enough stock available
			stock, err := l.svcCtx.StocksModel.FindOneBySkuIdWarehouseId(ctx, uint64(item.SkuId), uint64(in.WarehouseId))
			if err != nil {
				return err
			}
			if stock.Available < int64(item.Quantity) {
				return zeroerr.ErrInsufficientStock
			}

			// Verify and update stock
			err = l.svcCtx.StocksModel.DecrAvailable(ctx, uint64(item.SkuId), uint64(in.WarehouseId), int64(item.Quantity))
			if err != nil {
				return err
			}

			// Create stock record
			record := &model.StockRecords{
				SkuId:       uint64(item.SkuId),
				WarehouseId: uint64(in.WarehouseId),
				Type:        2, // Stock out
				Quantity:    int64(item.Quantity),
				OrderNo:     sql.NullString{String: in.OrderNo, Valid: len(in.OrderNo) > 0},
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
						Remark:        "out",
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

	return &inventory.CreateStockOutResponse{
		Success: true,
	}, nil
}
