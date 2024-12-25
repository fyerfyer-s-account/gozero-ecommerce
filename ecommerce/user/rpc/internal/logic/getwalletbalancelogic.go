package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWalletBalanceLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWalletBalanceLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWalletBalanceLogic {
	return &GetWalletBalanceLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 钱包操作
func (l *GetWalletBalanceLogic) GetWalletBalance(in *user.GetWalletBalanceRequest) (*user.GetWalletBalanceResponse, error) {
	// todo: add your logic here and delete this line

	return &user.GetWalletBalanceResponse{}, nil
}
