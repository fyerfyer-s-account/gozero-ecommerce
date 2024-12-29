package user

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProfileLogic {
	return &GetProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProfileLogic) GetProfile() (resp *types.UserInfo, err error) {
	// todo: add your logic here and delete this line
	// Get userId from JWT context
	userId := l.ctx.Value("userId").(int64)

	// Call user RPC
	userInfo, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
		UserId: userId,
	})
	if err != nil {
		logx.Errorf("get user profile error: %v", err)
		return nil, err
	}

	return &types.UserInfo{
		Id:          userInfo.UserInfo.UserId,
		Username:    userInfo.UserInfo.Username,
		Nickname:    userInfo.UserInfo.Nickname,
		Avatar:      userInfo.UserInfo.Avatar,
		Phone:       userInfo.UserInfo.Phone,
		Email:       userInfo.UserInfo.Email,
		Gender:      userInfo.UserInfo.Gender,
		MemberLevel: userInfo.UserInfo.MemberLevel,
		Balance:     userInfo.UserInfo.WalletBalance,
		CreatedAt:   userInfo.UserInfo.CreatedAt,
	}, nil
}
