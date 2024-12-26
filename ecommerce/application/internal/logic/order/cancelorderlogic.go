package order

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CancelOrderLogic) CancelOrder() error {
	// todo: add your logic here and delete this line

	return nil
}
