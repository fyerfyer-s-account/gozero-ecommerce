package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"

	"github.com/zeromicro/go-zero/core/logx"
)

type SelectItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSelectItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SelectItemLogic {
	return &SelectItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 商品选择
func (l *SelectItemLogic) SelectItem(in *cart.SelectItemRequest) (*cart.SelectItemResponse, error) {
	// todo: add your logic here and delete this line

	return &cart.SelectItemResponse{}, nil
}
