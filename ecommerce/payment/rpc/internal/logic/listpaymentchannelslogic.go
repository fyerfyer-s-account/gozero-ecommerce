package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListPaymentChannelsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListPaymentChannelsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListPaymentChannelsLogic {
	return &ListPaymentChannelsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListPaymentChannelsLogic) ListPaymentChannels(in *payment.ListPaymentChannelsRequest) (*payment.ListPaymentChannelsResponse, error) {
	// todo: add your logic here and delete this line

	return &payment.ListPaymentChannelsResponse{}, nil
}
