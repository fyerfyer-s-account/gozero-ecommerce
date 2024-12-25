package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReceiveOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveOrderLogic {
	return &ReceiveOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReceiveOrderLogic) ReceiveOrder(in *order.ReceiveOrderRequest) (*order.ReceiveOrderResponse, error) {
	// todo: add your logic here and delete this line

	return &order.ReceiveOrderResponse{}, nil
}
