package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &product.DeleteSkuResponse{}, nil
}
