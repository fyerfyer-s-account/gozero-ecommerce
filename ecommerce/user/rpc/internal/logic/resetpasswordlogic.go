package logic

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/cryptx"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 重置密码
func (l *ResetPasswordLogic) ResetPassword(in *user.ResetPasswordRequest) (*user.ResetPasswordResponse, error) {
	// 1. Find user by phone
	userInfo, err := l.svcCtx.UsersModel.FindOneByPhone(l.ctx, sql.NullString{String: in.Phone, Valid: true})
	if err != nil {
		if err == model.ErrNotFound {
			return nil, zeroerr.ErrPhoneNotFound
		}
		return nil, err
	}

	// 2. Verify code
	codeKey := fmt.Sprintf("reset:code:%s", in.Phone)
	code, err := l.svcCtx.BizRedis.Get(codeKey)
	if err != nil || code != in.VerifyCode {
		return nil, zeroerr.ErrInvalidVerifyCode
	}

	// 3. Validate new password
	if len(in.NewPassword) < l.svcCtx.Config.MinPasswordLength {
		return nil, zeroerr.ErrPasswordTooWeak
	}

	// 4. Update password
	newHashedPassword := cryptx.HashPassword(in.NewPassword, l.svcCtx.Config.Salt)
	err = l.svcCtx.UsersModel.UpdatePassword(l.ctx, userInfo.Id, newHashedPassword)
	if err != nil {
		logx.Errorf("reset password error: %v", err)
		return nil, zeroerr.ErrResetPasswordFailed
	}

	// 5. Clear verify code
	l.svcCtx.BizRedis.Del(codeKey)

	return &user.ResetPasswordResponse{
		Success: true,
	}, nil
}
