package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateSkuSalesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateSkuSalesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSkuSalesLogic {
	return &UpdateSkuSalesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateSkuSalesLogic) UpdateSkuSales(in *product.UpdateSkuSalesRequest) (*product.UpdateSkuSalesResponse, error) {
	if in.Id <= 0 {
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

	err = l.svcCtx.SkusModel.UpdateSales(l.ctx, uint64(in.Id), in.Increment)
	if err != nil {
		logx.Errorf("Failed to update SKU sales: %v", err)
		return nil, err
	}

	return &product.UpdateSkuSalesResponse{
		Success: true,
	}, nil
}
