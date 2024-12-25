package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateItemLogic {
	return &UpdateItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateItemLogic) UpdateItem(in *cart.UpdateItemRequest) (*cart.UpdateItemResponse, error) {
	// todo: add your logic here and delete this line

	return &cart.UpdateItemResponse{}, nil
}
