package payment

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	payment "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundNotifyLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRefundNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundNotifyLogic {
	return &RefundNotifyLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RefundNotifyLogic) RefundNotify(req *types.RefundNotifyReq) (resp *types.RefundNotifyResp, err error) {
	// todo: add your logic here and delete this line
	if req.PaymentType <= 0 || req.RefundNo == "" || req.Data == "" {
        return nil, zeroerr.ErrInvalidParameter
    }

    res, err := l.svcCtx.PaymentRpc.RefundNotify(l.ctx, &payment.RefundNotifyRequest{
        Channel:    int32(req.PaymentType),
        NotifyData: req.Data,
    })

    if err != nil {
        return &types.RefundNotifyResp{
            Code:    500,
            Message: err.Error(),
        }, nil
    }

    return &types.RefundNotifyResp{
        Code:    200,
        Message: res.ReturnMsg,
    }, nil
}
