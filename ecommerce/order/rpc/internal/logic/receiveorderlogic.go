package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

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
    if len(in.OrderNo) == 0 {
        return nil, zeroerr.ErrOrderNoEmpty
    }

    orderInfo, err := l.svcCtx.OrdersModel.FindOneByOrderNo(l.ctx, in.OrderNo)
    if err != nil {
        return nil, err
    }

    if orderInfo.Status != 3 { // 3: shipped
        return nil, zeroerr.ErrOrderStatusNotAllowed
    }

    // Update order status to completed
    err = l.svcCtx.OrdersModel.UpdateStatus(l.ctx, orderInfo.Id, 4)
    if err != nil {
        return nil, err
    }

    // Update shipping status
    shipping, err := l.svcCtx.OrderShippingModel.FindByOrderId(l.ctx, orderInfo.Id)
    if err != nil {
        return nil, err
    }

    shipping.Status = 2 // Received
    shipping.ReceiveTime = sql.NullTime {
		Time: time.Now(),
		Valid: true,
	}
    err = l.svcCtx.OrderShippingModel.Update(l.ctx, shipping)
    if err != nil {
        return nil, err
    }

    return &order.ReceiveOrderResponse{
        Success: true,
    }, nil
}
