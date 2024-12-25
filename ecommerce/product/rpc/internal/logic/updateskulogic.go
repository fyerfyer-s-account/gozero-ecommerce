package logic

import (
	"context"

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
	// todo: add your logic here and delete this line

	return &product.UpdateSkuResponse{}, nil
}
