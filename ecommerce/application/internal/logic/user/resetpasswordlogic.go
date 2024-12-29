package user

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ResetPasswordLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewResetPasswordLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ResetPasswordLogic {
	return &ResetPasswordLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ResetPasswordLogic) ResetPassword(req *types.ResetPasswordReq) error {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.UserRpc.ResetPassword(l.ctx, &user.ResetPasswordRequest{
		Phone:       req.Phone,
		VerifyCode:  req.Code,
		NewPassword: req.Password,
	})

	if err != nil {
		logx.Errorf("reset password error: %v", err)
		return err
	}

	return nil
}
