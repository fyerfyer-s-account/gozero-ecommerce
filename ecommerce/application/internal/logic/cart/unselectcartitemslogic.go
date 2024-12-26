package cart

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UnselectCartItemsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUnselectCartItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UnselectCartItemsLogic {
	return &UnselectCartItemsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UnselectCartItemsLogic) UnselectCartItems(req *types.BatchOperateReq) error {
	// todo: add your logic here and delete this line

	return nil
}
