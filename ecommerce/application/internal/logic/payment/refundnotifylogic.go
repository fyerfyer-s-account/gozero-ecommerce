package payment

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"

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

	return
}
