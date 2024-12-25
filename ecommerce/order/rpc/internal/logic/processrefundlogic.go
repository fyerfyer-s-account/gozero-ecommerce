package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessRefundLogic {
	return &ProcessRefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProcessRefundLogic) ProcessRefund(in *order.ProcessRefundRequest) (*order.ProcessRefundResponse, error) {
	// todo: add your logic here and delete this line

	return &order.ProcessRefundResponse{}, nil
}
