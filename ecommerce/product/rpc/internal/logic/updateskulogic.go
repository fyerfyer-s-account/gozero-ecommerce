package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSkuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSkuLogic {
	return &UpdateSkuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSkuLogic) UpdateSku(in *product.UpdateSkuRequest) (*product.UpdateSkuResponse, error) {
	// Validate input
	if in.Id <= 0 || in.Price <= 0 || in.Stock < 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check SKU exists
	_, err := l.svcCtx.SkusModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrSkuNotFound
	}

	// Update only price and stock
	err = l.svcCtx.SkusModel.UpdatePriceAndStock(l.ctx, uint64(in.Id), in.Price, in.Stock)
	if err != nil {
		logx.Errorf("Failed to update SKU: %v", err)
		return nil, zeroerr.ErrSkuUpdateFailed
	}

	return &product.UpdateSkuResponse{
		Success: true,
	}, nil
}
