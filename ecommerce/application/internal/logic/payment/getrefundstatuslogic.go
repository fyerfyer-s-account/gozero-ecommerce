package payment

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	payment "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRefundStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetRefundStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRefundStatusLogic {
	return &GetRefundStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetRefundStatusLogic) GetRefundStatus(req *types.RefundStatusReq) (resp *types.RefundStatusResp, err error) {
	// todo: add your logic here and delete this line

	if req.RefundNo == "" {
        return nil, zeroerr.ErrInvalidParameter
    }

    res, err := l.svcCtx.PaymentRpc.GetRefund(l.ctx, &payment.GetRefundRequest{
        RefundNo: req.RefundNo,
    })
    if err != nil {
        return nil, err
    }

    return &types.RefundStatusResp{
        Status:     int32(res.Refund.Status),
        Amount:     res.Refund.Amount,
        Reason:     res.Refund.Reason,
        RefundTime: res.Refund.RefundTime,
    }, nil
}
