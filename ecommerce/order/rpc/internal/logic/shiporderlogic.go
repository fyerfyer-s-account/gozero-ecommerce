package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShipOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShipOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShipOrderLogic {
	return &ShipOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ShipOrderLogic) ShipOrder(in *order.ShipOrderRequest) (*order.ShipOrderResponse, error) {
	// todo: add your logic here and delete this line

	return &order.ShipOrderResponse{}, nil
}
