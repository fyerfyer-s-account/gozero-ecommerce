package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
)

type BatchGetStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewBatchGetStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *BatchGetStockLogic {
	return &BatchGetStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *BatchGetStockLogic) BatchGetStock(in *inventory.BatchGetStockRequest) (*inventory.BatchGetStockResponse, error) {
	// Input validation
	if len(in.SkuIds) == 0 {
		return &inventory.BatchGetStockResponse{
			Stocks: make(map[int64]*inventory.Stock),
		}, nil
	}

	if in.WarehouseId <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Convert int64 slice to uint64 slice
	skuIds := make([]uint64, len(in.SkuIds))
	for i, id := range in.SkuIds {
		if id <= 0 {
			return nil, zeroerr.ErrInvalidParam
		}
		skuIds[i] = uint64(id)
	}

	// Query stocks from database
	stocks, err := l.svcCtx.StocksModel.BatchGet(l.ctx, skuIds, uint64(in.WarehouseId))
	if err != nil {
		l.Logger.Errorf("BatchGetStock error: %v", err)
		return nil, err
	}

	// Transform to response format
	result := make(map[int64]*inventory.Stock)
	for _, stock := range stocks {
		result[int64(stock.SkuId)] = &inventory.Stock{
			Id:          int64(stock.Id),
			SkuId:       int64(stock.SkuId),
			WarehouseId: int64(stock.WarehouseId),
			Available:   int32(stock.Available),
			Locked:      int32(stock.Locked),
			Total:       int32(stock.Total),
			CreatedAt:   stock.CreatedAt.Unix(),
			UpdatedAt:   stock.UpdatedAt.Unix(),
		}
	}

	return &inventory.BatchGetStockResponse{
		Stocks: result,
	}, nil
}
