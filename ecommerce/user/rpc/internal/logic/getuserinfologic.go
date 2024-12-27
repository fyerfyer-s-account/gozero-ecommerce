package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户信息
func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	// 1. Get user basic info
	userInfo, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		return nil, zeroerr.ErrUserNotFound
	}

	// 2. Get wallet info
	wallet, err := l.svcCtx.WalletAccountsModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil && err != model.ErrNotFound {
		return nil, err
	}

	return &user.GetUserInfoResponse{
		UserInfo: &user.UserInfo{
			UserId:        int64(userInfo.Id),
			Username:      userInfo.Username,
			Nickname:      userInfo.Nickname.String,
			Avatar:        userInfo.Avatar.String,
			Phone:         userInfo.Phone.String,
			Email:         userInfo.Email.String,
			Gender:        userInfo.Gender,
			WalletBalance: wallet.Balance,
			CreatedAt:     userInfo.CreatedAt.Unix(),
			UpdatedAt:     userInfo.UpdatedAt.Unix(),
		},
	}, nil
}
