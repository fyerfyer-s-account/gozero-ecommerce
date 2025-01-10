package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
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

func (l *GetPaymentStatusLogic) GetPaymentStatus(in *payment.GetPaymentStatusRequest) (*payment.GetPaymentStatusResponse, error) {
    if in.PaymentNo == "" {
        return nil, zeroerr.ErrInvalidParam
    }

    paymentOrder, err := l.svcCtx.PaymentOrdersModel.FindOneByPaymentNo(l.ctx, in.PaymentNo)
    if err != nil {
        return nil, zeroerr.ErrPaymentNotFound
    }

    var channelData string
    if paymentOrder.ChannelData.Valid {
        channelData = paymentOrder.ChannelData.String
    }

    return &payment.GetPaymentStatusResponse{
        Status:      int32(paymentOrder.Status),
        ChannelData: channelData,
    }, nil
}
