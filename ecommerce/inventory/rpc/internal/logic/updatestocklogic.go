package logic

import (
	"context"
	"database/sql"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
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
            // Stock update event
            updateEvent := &types.StockUpdatedEvent{
                InventoryEvent: types.InventoryEvent{
                    Type:        types.StockUpdated,
                    WarehouseID: int64(in.WarehouseId),
                    Timestamp:   record.CreatedAt,
                },
                SkuID:       int64(in.SkuId),
                OldQuantity: int32(stock.Available),
                NewQuantity: int32(stock.Available + int64(in.Quantity)),
                Reason:      in.Remark,
            }
            if err := l.svcCtx.Producer.PublishStockUpdated(ctx, updateEvent); err != nil {
                l.Logger.Error("Failed to publish stock update event", err)
            }

            // Check stock level and send alerts
            newQuantity := stock.Available + int64(in.Quantity)
            if newQuantity <= 0 {
                outOfStockEvent := &types.StockOutOfStockEvent{
                    InventoryEvent: types.InventoryEvent{
                        Type:        types.StockOutOfStock,
                        WarehouseID: int64(in.WarehouseId),
                        Timestamp:   record.CreatedAt,
                    },
                    SkuID:    int64(in.SkuId),
                    Quantity: 0,
                    Reason:   "Stock depleted after update",
                }
                if err := l.svcCtx.Producer.PublishStockOutOfStock(ctx, outOfStockEvent); err != nil {
                    l.Logger.Error("Failed to publish stock out event", err)
                }
            } else if newQuantity <= stock.AlertQuantity {
                lowStockEvent := &types.StockLowStockEvent{
                    InventoryEvent: types.InventoryEvent{
                        Type:        types.StockLowStock,
                        WarehouseID: int64(in.WarehouseId),
                        Timestamp:   record.CreatedAt,
                    },
                    SkuID:     int64(in.SkuId),
                    Quantity:  int32(newQuantity),
                    Threshold: int32(stock.AlertQuantity),
                }
                if err := l.svcCtx.Producer.PublishStockLowStock(ctx, lowStockEvent); err != nil {
                    l.Logger.Error("Failed to publish stock alert event", err)
                }
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
