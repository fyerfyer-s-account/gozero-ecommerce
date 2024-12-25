package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRefundLogic {
	return &CreateRefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 售后服务
func (l *CreateRefundLogic) CreateRefund(in *order.CreateRefundRequest) (*order.CreateRefundResponse, error) {
	// todo: add your logic here and delete this line

	return &order.CreateRefundResponse{}, nil
}
