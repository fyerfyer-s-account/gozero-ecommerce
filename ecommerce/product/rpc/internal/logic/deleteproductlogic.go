package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteProductLogic {
	return &DeleteProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteProductLogic) DeleteProduct(in *product.DeleteProductRequest) (*product.DeleteProductResponse, error) {
	// Validate input
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check product exists
	_, err := l.svcCtx.ProductsModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrProductNotFound
	}

	// Delete SKUs first
	err = l.svcCtx.SkusModel.DeleteByProductId(l.ctx, uint64(in.Id))
	if err != nil {
		logx.Errorf("Failed to delete SKUs for product %d: %v", in.Id, err)
		return nil, zeroerr.ErrProductDeleteFailed
	}

	// Delete product
	err = l.svcCtx.ProductsModel.Delete(l.ctx, uint64(in.Id))
	if err != nil {
		logx.Errorf("Failed to delete product %d: %v", in.Id, err)
		return nil, zeroerr.ErrProductDeleteFailed
	}

	return &product.DeleteProductResponse{
		Success: true,
	}, nil
}
