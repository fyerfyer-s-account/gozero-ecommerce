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

            oldQuantity := int32(stock.Available)

            // Verify and update stock
            err = l.svcCtx.StocksModel.DecrAvailable(ctx, uint64(item.SkuId), uint64(in.WarehouseId), int64(item.Quantity))
            if err != nil {
                return err
            }

            // Update stock totals
            stock.Available -= int64(item.Quantity)
            stock.Total -= int64(item.Quantity)

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
                // Stock update event
                updateEvent := &types.StockUpdatedEvent{
                    InventoryEvent: types.InventoryEvent{
                        Type:        types.StockUpdated,
                        WarehouseID: int64(in.WarehouseId),
                        Timestamp:   record.CreatedAt,
                    },
                    SkuID:       int64(item.SkuId),
                    OldQuantity: oldQuantity,
                    NewQuantity: int32(stock.Available),
                    Reason:      "stock_out",
                }
                if err := l.svcCtx.Producer.PublishStockUpdated(ctx, updateEvent); err != nil {
                    l.Logger.Error("Failed to publish stock update event", err)
                }

                // Check if stock is out
                if stock.Available == 0 {
                    outOfStockEvent := &types.StockOutOfStockEvent{
                        InventoryEvent: types.InventoryEvent{
                            Type:        types.StockOutOfStock,
                            WarehouseID: int64(in.WarehouseId),
                            Timestamp:   record.CreatedAt,
                        },
                        SkuID:    int64(item.SkuId),
                        Quantity: 0,
                        Reason:   "stock_depleted",
                    }
                    if err := l.svcCtx.Producer.PublishStockOutOfStock(ctx, outOfStockEvent); err != nil {
                        l.Logger.Error("Failed to publish stock out event", err)
                    }
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

    return &inventory.CreateStockOutResponse{
        Success: true,
    }, nil
}
