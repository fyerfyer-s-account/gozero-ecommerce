package payment

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"

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

func (l *GetRefundStatusLogic) GetRefundStatus() (resp *types.RefundStatusResp, err error) {
	// todo: add your logic here and delete this line

	return
}
