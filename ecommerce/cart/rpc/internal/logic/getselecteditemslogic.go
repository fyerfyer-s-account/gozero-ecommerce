package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSelectedItemsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSelectedItemsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSelectedItemsLogic {
	return &GetSelectedItemsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 结算相关
func (l *GetSelectedItemsLogic) GetSelectedItems(in *cart.GetSelectedItemsRequest) (*cart.GetSelectedItemsResponse, error) {
	// todo: add your logic here and delete this line

	return &cart.GetSelectedItemsResponse{}, nil
}
