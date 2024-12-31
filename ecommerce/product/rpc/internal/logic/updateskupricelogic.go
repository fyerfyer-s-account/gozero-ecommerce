package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateSkuPriceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSkuPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSkuPriceLogic {
	return &UpdateSkuPriceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSkuPriceLogic) UpdateSkuPrice(in *product.UpdateSkuPriceRequest) (*product.UpdateSkuPriceResponse, error) {
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	if in.Price <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check if SKU exists
	_, err := l.svcCtx.SkusModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, zeroerr.ErrSkuNotFound
		}
		return nil, err
	}

	err = l.svcCtx.SkusModel.UpdatePrice(l.ctx, uint64(in.Id), in.Price)
	if err != nil {
		logx.Errorf("Failed to update SKU price: %v", err)
		return nil, err
	}

	return &product.UpdateSkuPriceResponse{
		Success: true,
	}, nil
}
