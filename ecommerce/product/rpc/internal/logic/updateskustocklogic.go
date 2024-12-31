package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateSkuStockLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSkuStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSkuStockLogic {
	return &UpdateSkuStockLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSkuStockLogic) UpdateSkuStock(in *product.UpdateSkuStockRequest) (*product.UpdateSkuStockResponse, error) {
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check if SKU exists and get current stock
	sku, err := l.svcCtx.SkusModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, zeroerr.ErrSkuNotFound
		}
		return nil, err
	}

	// Prevent negative stock
	if sku.Stock+in.Increment < 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	err = l.svcCtx.SkusModel.UpdateStock(l.ctx, uint64(in.Id), in.Increment)
	if err != nil {
		logx.Errorf("Failed to update SKU stock: %v", err)
		return nil, err
	}

	return &product.UpdateSkuStockResponse{
		Success: true,
	}, nil
}
