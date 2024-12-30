package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteSkuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteSkuLogic {
	return &DeleteSkuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteSkuLogic) DeleteSku(in *product.DeleteSkuRequest) (*product.DeleteSkuResponse, error) {
	// Validate input
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check SKU exists
	_, err := l.svcCtx.SkusModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrSkuNotFound
	}

	// Delete SKU
	err = l.svcCtx.SkusModel.Delete(l.ctx, uint64(in.Id))
	if err != nil {
		logx.Errorf("Failed to delete SKU %d: %v", in.Id, err)
		return nil, zeroerr.ErrSkuDeleteFailed
	}

	return &product.DeleteSkuResponse{
		Success: true,
	}, nil
}
