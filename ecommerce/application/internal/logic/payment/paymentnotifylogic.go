package payment

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	payment "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentNotifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewPaymentNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentNotifyLogic {
	return &PaymentNotifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *PaymentNotifyLogic) PaymentNotify(req *types.PaymentNotifyReq) (resp *types.PaymentNotifyResp, err error) {
	// todo: add your logic here and delete this line
	if req.PaymentType == 0 || req.Data == "" {
        return nil, zeroerr.ErrInvalidParameter
    }

    res, err := l.svcCtx.PaymentRpc.PaymentNotify(l.ctx, &payment.PaymentNotifyRequest{
        Channel:    int32(req.PaymentType),
        NotifyData: req.Data,
    })
    if err != nil {
        return &types.PaymentNotifyResp{
            Code:    500,
            Message: err.Error(),
        }, nil
    }

    return &types.PaymentNotifyResp{
        Code:    200,
        Message: res.ReturnMsg,
    }, nil
}
