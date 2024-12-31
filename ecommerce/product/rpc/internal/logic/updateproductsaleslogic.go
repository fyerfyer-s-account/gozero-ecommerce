package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateProductSalesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductSalesLogic {
	return &UpdateProductSalesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductSalesLogic) UpdateProductSales(in *product.UpdateProductSalesRequest) (*product.UpdateProductSalesResponse, error) {
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	_, err := l.svcCtx.ProductsModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, zeroerr.ErrProductNotFound
		}
		return nil, err
	}

	err = l.svcCtx.ProductsModel.UpdateSales(l.ctx, uint64(in.Id), in.Increment)
	if err != nil {
		logx.Errorf("Failed to update product sales: %v", err)
		return nil, err
	}

	return &product.UpdateProductSalesResponse{
		Success: true,
	}, nil
}
