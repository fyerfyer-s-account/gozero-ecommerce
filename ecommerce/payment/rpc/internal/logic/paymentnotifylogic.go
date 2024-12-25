package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentNotifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPaymentNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentNotifyLogic {
	return &PaymentNotifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PaymentNotifyLogic) PaymentNotify(in *payment.PaymentNotifyRequest) (*payment.PaymentNotifyResponse, error) {
	// todo: add your logic here and delete this line

	return &payment.PaymentNotifyResponse{}, nil
}
