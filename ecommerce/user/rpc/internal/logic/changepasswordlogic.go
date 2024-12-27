package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/cryptx"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangePasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangePasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangePasswordLogic {
	return &ChangePasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 修改密码
func (l *ChangePasswordLogic) ChangePassword(in *user.ChangePasswordRequest) (*user.ChangePasswordResponse, error) {
	// 1. Get user info
	userInfo, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		if err == model.ErrNotFound {
			return nil, zeroerr.ErrUserNotFound
		}
		return nil, err
	}

	// 2. Verify old password
	oldHashedPassword := cryptx.HashPassword(in.OldPassword, l.svcCtx.Config.Salt)
	if oldHashedPassword != userInfo.Password {
		return nil, zeroerr.ErrOldPasswordIncorrect
	}

	// 3. Validate new password
	if len(in.NewPassword) < l.svcCtx.Config.MinPasswordLength {
		return nil, zeroerr.ErrPasswordTooWeak
	}

	// 4. Check if new password is same as old password
	if in.OldPassword == in.NewPassword {
		return nil, zeroerr.ErrSamePassword
	}

	// 5. Update password
	newHashedPassword := cryptx.HashPassword(in.NewPassword, l.svcCtx.Config.Salt)
	err = l.svcCtx.UsersModel.UpdatePassword(l.ctx, uint64(in.UserId), newHashedPassword)
	if err != nil {
		logx.Errorf("change password error: %v", err)
		return nil, zeroerr.ErrChangePasswordFailed
	}

	return &user.ChangePasswordResponse{
		Success: true,
	}, nil
}
