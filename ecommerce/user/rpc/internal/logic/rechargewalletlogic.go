package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type RechargeWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRechargeWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RechargeWalletLogic {
	return &RechargeWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RechargeWalletLogic) RechargeWallet(in *user.RechargeWalletRequest) (*user.RechargeWalletResponse, error) {
	// todo: add your logic here and delete this line

	return &user.RechargeWalletResponse{}, nil
}
