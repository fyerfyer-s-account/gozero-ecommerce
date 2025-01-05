package user

import (
	"context"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateProfileLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateProfileLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateProfileLogic {
	return &UpdateProfileLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateProfileLogic) UpdateProfile(req *types.UpdateProfileReq) error {
	// Get userId from JWT context with proper type conversion
	userIdVal := l.ctx.Value("userId")
	if userIdVal == nil {
		return zeroerr.ErrInvalidToken
	}

	// Handle json.Number type from JWT claims
	var userId int64
	var err error
	switch v := userIdVal.(type) {
	case json.Number:
		userId, err = v.Int64()
		if err != nil {
			logx.Errorf("invalid userId format: %v", err)
			return zeroerr.ErrInvalidToken
		}
	case float64:
		userId = int64(v)
	case int64:
		userId = v
	default:
		logx.Errorf("unexpected userId type: %T", userIdVal)
		return zeroerr.ErrInvalidToken
	}

	// Call RPC
	_, err = l.svcCtx.UserRpc.UpdateUserInfo(l.ctx, &user.UpdateUserInfoRequest{
		UserId:   userId,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Gender:   req.Gender,
		Phone:    req.Phone,
		Email:    req.Email,
	})

	if err != nil {
		logx.Errorf("update profile error: %v", err)
		return err
	}

	return nil
}
