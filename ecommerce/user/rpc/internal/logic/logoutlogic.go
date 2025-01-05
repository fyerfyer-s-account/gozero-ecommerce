package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/jwtx"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LogoutLogic) Logout(in *user.LogoutRequest) (*user.LogoutResponse, error) {
	claims, err := jwtx.ParseToken(in.AccessToken, l.svcCtx.Config.JwtAuth.AccessSecret)
	if err != nil {
		return nil, zeroerr.ErrInvalidToken
	}

	// Set user offline
	err = l.svcCtx.UsersModel.UpdateOnline(l.ctx, uint64(claims.UserID), 0)
	if err != nil {
		return nil, err
	}

	return &user.LogoutResponse{
		Success: true,
	}, nil
}
