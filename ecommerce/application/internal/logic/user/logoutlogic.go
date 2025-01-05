package user

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LogoutLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewLogoutLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LogoutLogic {
	return &LogoutLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LogoutLogic) Logout(req *types.LogoutReq) (resp *types.LogoutResp, err error) {
    // Call RPC to set user offline
    _, err = l.svcCtx.UserRpc.Logout(l.ctx, &user.LogoutRequest{
        AccessToken: req.AccessToken,
    })
    if err != nil {
        return nil, err
    }

    // Add token to blacklist
    err = l.svcCtx.TokenBlacklist.SetexCtx(
        l.ctx,
        req.AccessToken,
        "1",
        int(l.svcCtx.Config.Auth.AccessExpire),
    )
    if err != nil {
        logx.Error("failed to blacklist token:", err)
    }

    return &types.LogoutResp{
        Success: true,
    }, nil
}
