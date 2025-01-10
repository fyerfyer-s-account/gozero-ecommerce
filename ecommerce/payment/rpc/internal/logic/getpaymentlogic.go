package logic

import (
	"context"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentLogic {
	return &GetPaymentLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPaymentLogic) GetPayment(in *payment.GetPaymentRequest) (*payment.GetPaymentResponse, error) {
	if in.PaymentNo == "" {
		return nil, zeroerr.ErrInvalidParam
	}

	paymentOrder, err := l.svcCtx.PaymentOrdersModel.FindOneByPaymentNo(l.ctx, in.PaymentNo)
	if err != nil {
		return nil, zeroerr.ErrPaymentNotFound
	}

	var channelData, notifyUrl, returnUrl string
	if paymentOrder.ChannelData.Valid {
		channelData = paymentOrder.ChannelData.String
	}
	if paymentOrder.NotifyUrl.Valid {
		notifyUrl = paymentOrder.NotifyUrl.String
	}
	if paymentOrder.ReturnUrl.Valid {
		returnUrl = paymentOrder.ReturnUrl.String
	}

	resp := &payment.GetPaymentResponse{
		Payment: &payment.PaymentOrder{
			PaymentNo:   paymentOrder.PaymentNo,
			OrderNo:     paymentOrder.OrderNo,
			UserId:      int64(paymentOrder.UserId),
			Amount:      paymentOrder.Amount,
			Channel:     paymentOrder.Channel,
			ChannelData: channelData,
			Status:      int32(paymentOrder.Status),
			NotifyUrl:   notifyUrl,
			ReturnUrl:   returnUrl,
			CreatedAt:   paymentOrder.CreatedAt.Unix(),
			UpdatedAt:   paymentOrder.UpdatedAt.Unix(),
		},
	}

	if paymentOrder.ExpireTime.Valid {
		resp.Payment.ExpireTime = paymentOrder.ExpireTime.Time.Unix()
	}
	if paymentOrder.PayTime.Valid {
		resp.Payment.PayTime = paymentOrder.PayTime.Time.Unix()
	}

	return resp, nil
}
