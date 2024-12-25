// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: payment.proto

package server

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/logic"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
)

type PaymentServer struct {
	svcCtx *svc.ServiceContext
	payment.UnimplementedPaymentServer
}

func NewPaymentServer(svcCtx *svc.ServiceContext) *PaymentServer {
	return &PaymentServer{
		svcCtx: svcCtx,
	}
}

// 支付相关
func (s *PaymentServer) CreatePayment(ctx context.Context, in *payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
	l := logic.NewCreatePaymentLogic(ctx, s.svcCtx)
	return l.CreatePayment(in)
}

func (s *PaymentServer) GetPayment(ctx context.Context, in *payment.GetPaymentRequest) (*payment.GetPaymentResponse, error) {
	l := logic.NewGetPaymentLogic(ctx, s.svcCtx)
	return l.GetPayment(in)
}

func (s *PaymentServer) PaymentNotify(ctx context.Context, in *payment.PaymentNotifyRequest) (*payment.PaymentNotifyResponse, error) {
	l := logic.NewPaymentNotifyLogic(ctx, s.svcCtx)
	return l.PaymentNotify(in)
}

func (s *PaymentServer) GetPaymentStatus(ctx context.Context, in *payment.GetPaymentStatusRequest) (*payment.GetPaymentStatusResponse, error) {
	l := logic.NewGetPaymentStatusLogic(ctx, s.svcCtx)
	return l.GetPaymentStatus(in)
}

// 退款相关
func (s *PaymentServer) CreateRefund(ctx context.Context, in *payment.CreateRefundRequest) (*payment.CreateRefundResponse, error) {
	l := logic.NewCreateRefundLogic(ctx, s.svcCtx)
	return l.CreateRefund(in)
}

func (s *PaymentServer) GetRefund(ctx context.Context, in *payment.GetRefundRequest) (*payment.GetRefundResponse, error) {
	l := logic.NewGetRefundLogic(ctx, s.svcCtx)
	return l.GetRefund(in)
}

func (s *PaymentServer) RefundNotify(ctx context.Context, in *payment.RefundNotifyRequest) (*payment.RefundNotifyResponse, error) {
	l := logic.NewRefundNotifyLogic(ctx, s.svcCtx)
	return l.RefundNotify(in)
}

// 支付渠道
func (s *PaymentServer) CreatePaymentChannel(ctx context.Context, in *payment.CreatePaymentChannelRequest) (*payment.CreatePaymentChannelResponse, error) {
	l := logic.NewCreatePaymentChannelLogic(ctx, s.svcCtx)
	return l.CreatePaymentChannel(in)
}

func (s *PaymentServer) UpdatePaymentChannel(ctx context.Context, in *payment.UpdatePaymentChannelRequest) (*payment.UpdatePaymentChannelResponse, error) {
	l := logic.NewUpdatePaymentChannelLogic(ctx, s.svcCtx)
	return l.UpdatePaymentChannel(in)
}

func (s *PaymentServer) ListPaymentChannels(ctx context.Context, in *payment.ListPaymentChannelsRequest) (*payment.ListPaymentChannelsResponse, error) {
	l := logic.NewListPaymentChannelsLogic(ctx, s.svcCtx)
	return l.ListPaymentChannels(in)
}
