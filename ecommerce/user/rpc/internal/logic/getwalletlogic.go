package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetWalletLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetWalletLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetWalletLogic {
	return &GetWalletLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 钱包操作
func (l *GetWalletLogic) GetWallet(in *user.GetWalletRequest) (*user.GetWalletResponse, error) {
	// todo: add your logic here and delete this line
	// 1. Check if user exists
	_, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, zeroerr.ErrUserNotFound
		}
		return nil, err
	}

	// 2. Get wallet info
	wallet, err := l.svcCtx.WalletAccountsModel.FindOneByUserId(l.ctx, uint64(in.UserId))
	if err != nil {
		if err == model.ErrNotFound {
			// Return empty wallet for new users
			return &user.GetWalletResponse{
				Balance:      0,
				Status:       1, // Normal status
				FreezeAmount: 0,
			}, nil
		}
		logx.Errorf("get wallet error: %v", err)
		return nil, zeroerr.ErrInvalidAmount
	}

	// 3. Check wallet status
	if wallet.Status != 1 {
		return nil, zeroerr.ErrWalletDisabled
	}

	return &user.GetWalletResponse{
		Balance:      wallet.Balance,
		Status:       wallet.Status,
		FreezeAmount: wallet.FrozenAmount,
	}, nil
}
