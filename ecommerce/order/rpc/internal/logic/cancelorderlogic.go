package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelOrderLogic) CancelOrder(in *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
    if len(in.OrderNo) == 0 {
        return nil, zeroerr.ErrOrderInvalidParam
    }

    orderInfo, err := l.svcCtx.OrdersModel.FindOneByOrderNo(l.ctx, in.OrderNo)
    if err != nil {
        return nil, err
    }

    if orderInfo.Status != 1 {
        return nil, zeroerr.ErrOrderStatusInvalid
    }

    err = l.svcCtx.OrdersModel.UpdateStatus(l.ctx, orderInfo.Id, 5) // 5: Canceled
    if err != nil {
        return nil, err
    }

    return &order.CancelOrderResponse{
        Success: true,
    }, nil
}
