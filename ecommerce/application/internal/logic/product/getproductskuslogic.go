package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductSkusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductSkusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductSkusLogic {
	return &GetProductSkusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductSkusLogic) GetProductSkus() (resp []types.Sku, err error) {
	// todo: add your logic here and delete this line

	return
}
