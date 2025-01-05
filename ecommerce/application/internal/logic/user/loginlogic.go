package user

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.TokenResp, err error) {
	res, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginRequest{
		Username: req.Username,
		Password: req.Password,
	})

	if err != nil {
		logx.Errorf("login error: %v", err)
		return nil, err
	}

	// Make sure to use the same Auth.AccessSecret as middleware
	if res.AccessToken == "" {
		logx.Error("empty token received from user RPC")
		return nil, zeroerr.ErrInvalidToken
	}

	return &types.TokenResp{
		AccessToken:  res.AccessToken,
		RefreshToken: res.RefreshToken,
		ExpiresIn:    res.ExpiresIn,
	}, nil
}
