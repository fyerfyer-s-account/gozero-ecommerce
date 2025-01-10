package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRefundLogic {
	return &GetRefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRefundLogic) GetRefund(in *payment.GetRefundRequest) (*payment.GetRefundResponse, error) {
    if in.RefundNo == "" {
        return nil, zeroerr.ErrInvalidParam
    }

    refundOrder, err := l.svcCtx.RefundOrdersModel.FindOneByRefundNo(l.ctx, in.RefundNo)
    if err != nil {
        return nil, zeroerr.ErrRefundNotFound
    }

    var notifyUrl string
    if refundOrder.NotifyUrl.Valid {
        notifyUrl = refundOrder.NotifyUrl.String
    }

    resp := &payment.GetRefundResponse{
        Refund: &payment.RefundOrder{
            RefundNo:   refundOrder.RefundNo,
            PaymentNo:  refundOrder.PaymentNo,
            OrderNo:    refundOrder.OrderNo,
            UserId:     int64(refundOrder.UserId),
            Amount:     refundOrder.Amount,
            Reason:     refundOrder.Reason,
            Status:     int32(refundOrder.Status),
            NotifyUrl:  notifyUrl,
            CreatedAt:  refundOrder.CreatedAt.Unix(),
            UpdatedAt:  refundOrder.UpdatedAt.Unix(),
        },
    }

    if refundOrder.RefundTime.Valid {
        resp.Refund.RefundTime = refundOrder.RefundTime.Time.Unix()
    }

    return resp, nil
}
