package cart

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectCartItemsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewSelectCartItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectCartItemsLogic {
	return &SelectCartItemsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *SelectCartItemsLogic) SelectCartItems(req *types.BatchOperateReq) error {
	// todo: add your logic here and delete this line

	return nil
}
