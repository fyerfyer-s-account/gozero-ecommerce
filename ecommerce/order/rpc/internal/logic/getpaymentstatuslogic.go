package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetPaymentStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetPaymentStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetPaymentStatusLogic {
	return &GetPaymentStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetPaymentStatusLogic) GetPaymentStatus(in *order.GetPaymentStatusRequest) (*order.GetPaymentStatusResponse, error) {
	// todo: add your logic here and delete this line

	return &order.GetPaymentStatusResponse{}, nil
}
