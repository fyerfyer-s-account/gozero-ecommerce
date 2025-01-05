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
	// Get userId from JWT context with proper type conversion
	userIdVal := l.ctx.Value("userId")
	if userIdVal == nil {
		return nil, zeroerr.ErrInvalidToken
	}

	// Handle json.Number type from JWT claims
	var userId int64
	switch v := userIdVal.(type) {
	case json.Number:
		userId, err = v.Int64()
		if err != nil {
			logx.Errorf("invalid userId format: %v", err)
			return nil, zeroerr.ErrInvalidToken
		}
	case float64:
		userId = int64(v)
	case int64:
		userId = v
	default:
		logx.Errorf("unexpected userId type: %T", userIdVal)
		return nil, zeroerr.ErrInvalidToken
	}

	// Call user RPC with error handling
	userInfo, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user.GetUserInfoRequest{
		UserId: userId,
	})
	if err != nil {
		logx.Errorf("get user profile error: %v", err)
		return nil, zeroerr.ErrUserNotFound
	}

	if userInfo == nil || userInfo.UserInfo == nil {
		return nil, zeroerr.ErrUserNotFound
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
