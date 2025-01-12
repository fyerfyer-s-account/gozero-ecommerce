package order

import (
	"context"
	"strconv"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	order "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelOrderLogic) CancelOrder(req *types.CancelOrderReq) error {
    orderId := strconv.FormatInt(req.Id, 10)

    // Get order first to get order_no
    orderResp, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderRequest{
        OrderNo: orderId,
    })
    if err != nil {
        return err
    }

    // Call cancel order RPC
    _, err = l.svcCtx.OrderRpc.CancelOrder(l.ctx, &order.CancelOrderRequest{
        OrderNo: orderResp.Order.OrderNo,
        Reason:  "用户取消",
    })

    return err
}
