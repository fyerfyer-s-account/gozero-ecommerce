package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentChannelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreatePaymentChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentChannelLogic {
	return &CreatePaymentChannelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 支付渠道
func (l *CreatePaymentChannelLogic) CreatePaymentChannel(in *payment.CreatePaymentChannelRequest) (*payment.CreatePaymentChannelResponse, error) {
	// todo: add your logic here and delete this line

	return &payment.CreatePaymentChannelResponse{}, nil
}
