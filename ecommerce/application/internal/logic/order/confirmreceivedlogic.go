package order

import (
	"context"
	"strconv"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	order "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmReceivedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmReceivedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmReceivedLogic {
	return &ConfirmReceivedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmReceivedLogic) ConfirmReceived(req *types.ConfirmOrderReq) error {
    orderId := strconv.FormatInt(req.Id, 10)

    // Get order details first
    orderResp, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderRequest{
        OrderNo: orderId,
    })
    if err != nil {
        return err
    }
	
    _, err = l.svcCtx.OrderRpc.ReceiveOrder(l.ctx, &order.ReceiveOrderRequest{
        OrderNo: orderResp.Order.OrderNo,
    })

    return err
}
