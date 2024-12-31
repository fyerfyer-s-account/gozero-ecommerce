package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateProductPriceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductPriceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductPriceLogic {
	return &UpdateProductPriceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductPriceLogic) UpdateProductPrice(in *product.UpdateProductPriceRequest) (*product.UpdateProductPriceResponse, error) {
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	if in.Price <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	_, err := l.svcCtx.ProductsModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, zeroerr.ErrProductNotFound
		}
		return nil, err
	}

	err = l.svcCtx.ProductsModel.UpdatePrice(l.ctx, uint64(in.Id), in.Price)
	if err != nil {
		logx.Errorf("Failed to update product price: %v", err)
		return nil, zeroerr.ErrProductUpdateFailed
	}

	return &product.UpdateProductPriceResponse{
		Success: true,
	}, nil
}
