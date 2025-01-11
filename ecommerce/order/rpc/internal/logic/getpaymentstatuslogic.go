package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentStatusLogic {
	return &GetPaymentStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPaymentStatusLogic) GetPaymentStatus(in *order.GetPaymentStatusRequest) (*order.GetPaymentStatusResponse, error) {
    if len(in.PaymentNo) == 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    payment, err := l.svcCtx.OrderPaymentsModel.FindOneByPaymentNo(l.ctx, in.PaymentNo)
    if err != nil {
        return nil, err
    }

    return &order.GetPaymentStatusResponse{
        Status: int32(payment.Status),
    }, nil
}
