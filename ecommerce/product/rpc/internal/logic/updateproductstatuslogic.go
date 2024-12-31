package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type UpdateProductStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateProductStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProductStatusLogic {
	return &UpdateProductStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateProductStatusLogic) UpdateProductStatus(in *product.UpdateProductStatusRequest) (*product.UpdateProductStatusResponse, error) {
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	if in.Status != 1 && in.Status != 2 { // 1: 上架, 2: 下架
		return nil, zeroerr.ErrInvalidParam
	}

	_, err := l.svcCtx.ProductsModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		if err == sqlx.ErrNotFound {
			return nil, zeroerr.ErrProductNotFound
		}
		return nil, err
	}

	err = l.svcCtx.ProductsModel.UpdateStatus(l.ctx, uint64(in.Id), in.Status)
	if err != nil {
		logx.Errorf("Failed to update product status: %v", err)
		return nil, err
	}

	return &product.UpdateProductStatusResponse{
		Success: true,
	}, nil
}
