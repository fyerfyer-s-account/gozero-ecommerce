package logic

import (
	"context"
	"database/sql"

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

	// 2. Validate gender value
	if in.Gender != "" && in.Gender != "M" && in.Gender != "F" {
		return nil, zeroerr.ErrInvalidGender
	}

	// 3. Update user info
	err = l.svcCtx.UsersModel.Update(l.ctx, &model.Users{
		Id:       uint64(in.UserId),
		Nickname: sql.NullString{String: in.Nickname, Valid: len(in.Nickname) > 0},
		Avatar:   sql.NullString{String: in.Avatar, Valid: len(in.Avatar) > 0},
		Gender:   in.Gender,
		Phone:    sql.NullString{String: in.Phone, Valid: len(in.Phone) > 0},
		Email:    sql.NullString{String: in.Email, Valid: len(in.Email) > 0},
	})

	if err != nil {
		logx.Errorf("update user info error: %v", err)
		return nil, zeroerr.ErrUpdateProfile
	}

	return &user.UpdateUserInfoResponse{
		Success: true,
	}, nil
}
