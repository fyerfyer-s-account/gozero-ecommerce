package cart

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSelectedItemsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetSelectedItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelectedItemsLogic {
	return &GetSelectedItemsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetSelectedItemsLogic) GetSelectedItems() (resp *types.SelectedItemsResp, err error) {
	// todo: add your logic here and delete this line

	return
}
