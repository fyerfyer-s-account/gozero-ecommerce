package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type WithdrawWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewWithdrawWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *WithdrawWalletLogic {
	return &WithdrawWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *WithdrawWalletLogic) WithdrawWallet(in *user.WithdrawWalletRequest) (*user.WithdrawWalletResponse, error) {
	// todo: add your logic here and delete this line

	return &user.WithdrawWalletResponse{}, nil
}
