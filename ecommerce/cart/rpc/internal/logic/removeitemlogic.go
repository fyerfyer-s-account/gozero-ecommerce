package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type RemoveItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRemoveItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RemoveItemLogic {
	return &RemoveItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RemoveItemLogic) RemoveItem(in *cart.RemoveItemRequest) (*cart.RemoveItemResponse, error) {
	// todo: add your logic here and delete this line

	return &cart.RemoveItemResponse{}, nil
}
