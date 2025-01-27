package logic

import (
	"context"
	"database/sql"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
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

            var oldQuantity int32
            if err == model.ErrNotFound {
                // Create new stock
                stock = &model.Stocks{
                    SkuId:         uint64(item.SkuId),
                    WarehouseId:   uint64(in.WarehouseId),
                    Available:     int64(item.Quantity),
                    Locked:        0,
                    Total:         int64(item.Quantity),
                    AlertQuantity: 10, // Default alert threshold
                }
                _, err = l.svcCtx.StocksModel.Insert(ctx, stock)
                if err != nil {
                    return err
                }
            } else {
                oldQuantity = int32(stock.Available)
                // Update existing stock
                err = l.svcCtx.StocksModel.IncrAvailable(ctx, uint64(item.SkuId), uint64(in.WarehouseId), int64(item.Quantity))
                if err != nil {
                    return err
                }
                stock.Available += int64(item.Quantity)
                stock.Total += int64(item.Quantity)
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
                event := &types.StockUpdatedEvent{
                    InventoryEvent: types.InventoryEvent{
                        Type:        types.StockUpdated,
                        WarehouseID: int64(in.WarehouseId),
                        Timestamp:   record.CreatedAt,
                    },
                    SkuID:       int64(item.SkuId),
                    OldQuantity: oldQuantity,
                    NewQuantity: int32(stock.Available),
                    Reason:      "stock_in",
                }
                if err := l.svcCtx.Producer.PublishStockUpdated(ctx, event); err != nil {
                    l.Logger.Error("Failed to publish stock update event", err)
                }

                // Check for low stock alert
                if stock.Available <= stock.AlertQuantity {
                    alertEvent := &types.StockLowStockEvent{
                        InventoryEvent: types.InventoryEvent{
                            Type:        types.StockLowStock,
                            WarehouseID: int64(in.WarehouseId),
                            Timestamp:   record.CreatedAt,
                        },
                        SkuID:     int64(item.SkuId),
                        Quantity:  int32(stock.Available),
                        Threshold: int32(stock.AlertQuantity),
                    }
                    if err := l.svcCtx.Producer.PublishStockLowStock(ctx, alertEvent); err != nil {
                        l.Logger.Error("Failed to publish stock alert event", err)
                    }
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