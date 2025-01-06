package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateSkuLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateSkuLogic {
	return &UpdateSkuLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateSkuLogic) UpdateSku(req *types.UpdateSkuReq) error {
	// todo: add your logic here and delete this line

	return nil
}
