package user

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type RechargeWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewRechargeWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RechargeWalletLogic {
	return &RechargeWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RechargeWalletLogic) RechargeWallet(req *types.RechargeReq) error {
	// todo: add your logic here and delete this line

	return nil
}
