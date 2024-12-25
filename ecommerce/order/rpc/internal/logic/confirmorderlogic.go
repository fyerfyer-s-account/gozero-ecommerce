package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewConfirmOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmOrderLogic {
	return &ConfirmOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 订单履约
func (l *ConfirmOrderLogic) ConfirmOrder(in *order.ConfirmOrderRequest) (*order.ConfirmOrderResponse, error) {
	// todo: add your logic here and delete this line

	return &order.ConfirmOrderResponse{}, nil
}
