package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundNotifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundNotifyLogic {
	return &RefundNotifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundNotifyLogic) RefundNotify(in *payment.RefundNotifyRequest) (*payment.RefundNotifyResponse, error) {
	// todo: add your logic here and delete this line

	return &payment.RefundNotifyResponse{}, nil
}
