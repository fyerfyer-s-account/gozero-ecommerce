package payment

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	payment "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetPaymentStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentStatusLogic {
	return &GetPaymentStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetPaymentStatusLogic) GetPaymentStatus(req *types.PaymentStatusReq) (resp *types.PaymentStatusResp, err error) {
	// todo: add your logic here and delete this line
	if req.PaymentNo == "" {
        return nil, zeroerr.ErrPaymentNoEmpty
    }

    res, err := l.svcCtx.PaymentRpc.GetPayment(l.ctx, &payment.GetPaymentRequest{
        PaymentNo: req.PaymentNo,
    })
    if err != nil {
        return nil, err
    }

    return &types.PaymentStatusResp{
        Status:   int32(res.Payment.Status),
        Amount:   res.Payment.Amount,
        PayTime:  res.Payment.PayTime,
        ErrorMsg: "",
    }, nil
}
