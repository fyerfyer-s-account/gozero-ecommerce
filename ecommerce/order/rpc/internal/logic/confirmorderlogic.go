package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

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

func (l *ConfirmOrderLogic) ConfirmOrder(in *order.ConfirmOrderRequest) (*order.ConfirmOrderResponse, error) {
    if len(in.OrderNo) == 0 {
        return nil, zeroerr.ErrOrderNoEmpty
    }

    orderInfo, err := l.svcCtx.OrdersModel.FindOneByOrderNo(l.ctx, in.OrderNo)
    if err != nil {
        return nil, err
    }

    if orderInfo.Status != 2 { // 2: paid, waiting for shipment
        return nil, zeroerr.ErrOrderStatusNotAllowed
    }

    err = l.svcCtx.OrdersModel.UpdateStatus(l.ctx, orderInfo.Id, 3) // 3: shipped
    if err != nil {
        return nil, err
    }

    return &order.ConfirmOrderResponse{
        Success: true,
    }, nil
}
