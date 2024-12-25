package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSkuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSkuLogic {
	return &GetSkuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSkuLogic) GetSku(in *product.GetSkuRequest) (*product.GetSkuResponse, error) {
	// todo: add your logic here and delete this line

	return &product.GetSkuResponse{}, nil
}
