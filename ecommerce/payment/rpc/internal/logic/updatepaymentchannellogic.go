package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePaymentChannelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePaymentChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePaymentChannelLogic {
	return &UpdatePaymentChannelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePaymentChannelLogic) UpdatePaymentChannel(in *payment.UpdatePaymentChannelRequest) (*payment.UpdatePaymentChannelResponse, error) {
	// todo: add your logic here and delete this line

	return &payment.UpdatePaymentChannelResponse{}, nil
}
