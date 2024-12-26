package order

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/zeromicro/go-zero/core/logx"
)

type ConfirmReceivedLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewConfirmReceivedLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfirmReceivedLogic {
	return &ConfirmReceivedLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfirmReceivedLogic) ConfirmReceived() error {
	// todo: add your logic here and delete this line

	return nil
}
