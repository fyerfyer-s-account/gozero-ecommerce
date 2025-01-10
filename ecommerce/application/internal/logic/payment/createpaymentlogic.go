package payment

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	payment "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
	return &CreatePaymentLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreatePaymentLogic) CreatePayment(req *types.CreatePaymentReq) (resp *types.CreatePaymentResp, err error) {
	// todo: add your logic here and delete this line
	userId := l.ctx.Value("userId").(int64)

    res, err := l.svcCtx.PaymentRpc.CreatePayment(l.ctx, &payment.CreatePaymentRequest{
        OrderNo:   req.OrderNo,
        UserId:    userId,
        Amount:    req.Amount,
        Channel:   int64(req.PaymentType),
        NotifyUrl: req.NotifyUrl,
        ReturnUrl: req.ReturnUrl,
    })

    if err != nil {
        return nil, err
    }

    return &types.CreatePaymentResp{
        PaymentNo: res.PaymentNo,
        PayUrl:    res.PayUrl,
    }, nil
}
