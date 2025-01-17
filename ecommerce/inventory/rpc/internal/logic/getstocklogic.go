package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetStockLogic {
	return &GetStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 库存管理
func (l *GetStockLogic) GetStock(in *inventory.GetStockRequest) (*inventory.GetStockResponse, error) {
    // Input validation
    if in.SkuId <= 0 || in.WarehouseId <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Query stock from database
    stock, err := l.svcCtx.StocksModel.FindOneBySkuIdWarehouseId(l.ctx, uint64(in.SkuId), uint64(in.WarehouseId))
    if err != nil {
        if err == model.ErrNotFound {
            return nil, zeroerr.ErrStockNotFound
        }
        return nil, err
    }

    // Transform to response
    return &inventory.GetStockResponse{
        Stock: &inventory.Stock{
            Id:          int64(stock.Id),
            SkuId:       int64(stock.SkuId),
            WarehouseId: int64(stock.WarehouseId),
            Available:   int32(stock.Available),
            Locked:      int32(stock.Locked),
            Total:       int32(stock.Total),
            CreatedAt:   stock.CreatedAt.Unix(),
            UpdatedAt:   stock.UpdatedAt.Unix(),
        },
    }, nil
}
