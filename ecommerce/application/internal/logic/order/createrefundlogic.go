package order

import (
	"context"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	order "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRefundLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRefundLogic {
	return &CreateRefundLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateRefundLogic) CreateRefund(req *types.CreateRefundReq) (*types.RefundInfo, error) {
    // Get order details first
    orderResp, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderRequest{
        OrderNo: req.OrderNo,
    })
    if err != nil {
        return nil, err
    }

    if orderResp.Order.PayAmount < req.Amount {
        return nil, zeroerr.ErrRefundExceedAmount
    }

    // Create refund via RPC
    refundResp, err := l.svcCtx.OrderRpc.CreateRefund(l.ctx, &order.CreateRefundRequest{
        OrderNo:     orderResp.Order.OrderNo,
        Amount:      req.Amount,
        Reason:      req.Reason,
        Description: req.Desc,
        Images:      req.Images,
    })
    if err != nil {
        return nil, err
    }

    return &types.RefundInfo{
		Id:        req.Id,
        RefundNo:  refundResp.RefundNo,
        Status:    0, // Initial status: pending
        Amount:    req.Amount,
        Reason:    req.Reason,
        Desc:      req.Desc,
        Images:    req.Images,
        CreatedAt: time.Now().Unix(),
    }, nil
}
