package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateUserInfoLogic {
	return &UpdateUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 更新用户信息
func (l *UpdateUserInfoLogic) UpdateUserInfo(in *user.UpdateUserInfoRequest) (*user.UpdateUserInfoResponse, error) {
	// 1. Check if user exists
	_, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, zeroerr.ErrUserNotFound
		}
		return nil, err
	}

	// 3. Update user info
	err = l.svcCtx.UsersModel.UpdateProfile(l.ctx, uint64(in.UserId), in.Nickname, in.Gender)

	if err != nil {
		logx.Errorf("update user info error: %v", err)
		return nil, zeroerr.ErrUpdateProfile
	}

	return &user.UpdateUserInfoResponse{
		Success: true,
	}, nil
}
