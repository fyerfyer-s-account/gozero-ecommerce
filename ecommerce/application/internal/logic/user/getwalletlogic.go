package user

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWalletLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWalletLogic {
	return &GetWalletLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetWalletLogic) GetWallet() (resp *types.WalletDetail, err error) {
	// todo: add your logi// Get userId from JWT context
	userId := l.ctx.Value("userId").(int64)

	// Call user RPC
	wallet, err := l.svcCtx.UserRpc.GetWallet(l.ctx, &user.GetWalletRequest{
		UserId: userId,
	})
	if err != nil {
		logx.Errorf("get wallet error: %v", err)
		return nil, err
	}

	return &types.WalletDetail{
		Balance:      wallet.Balance,
		Status:       wallet.Status,
		FrozenAmount: wallet.FreezeAmount,
	}, nil
}
