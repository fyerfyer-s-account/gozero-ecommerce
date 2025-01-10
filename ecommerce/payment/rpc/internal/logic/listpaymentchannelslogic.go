package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
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
    var channels []*model.PaymentChannels
    var err error

    if in.Status > 0 {
        channels, err = l.svcCtx.PaymentChannelsModel.FindManyByStatus(l.ctx, int64(in.Status))
    } else {
        channels, err = l.svcCtx.PaymentChannelsModel.FindAll(l.ctx)
    }

    if err != nil {
        return nil, err
    }

    var protoChannels []*payment.PaymentChannel
    for _, ch := range channels {
        protoChannels = append(protoChannels, &payment.PaymentChannel{
            Id:        int64(ch.Id),
            Name:      ch.Name,
            Channel:   ch.Channel,
            Config:    ch.Config,
            Status:    int32(ch.Status),
            CreatedAt: ch.CreatedAt.Unix(),
            UpdatedAt: ch.UpdatedAt.Unix(),
        })
    }

    return &payment.ListPaymentChannelsResponse{
        Channels: protoChannels,
    }, nil
}
