package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/cryptx"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/jwtx"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 用户登录
func (l *LoginLogic) Login(in *user.LoginRequest) (*user.LoginResponse, error) {
	// 1. Validate input
	if len(in.Username) == 0 || len(in.Password) == 0 {
		return nil, zeroerr.ErrInvalidUsername
	}

	// 2. Find user
	userInfo, err := l.svcCtx.UsersModel.FindOneByPhoneOrEmail(l.ctx, in.Username)
	if err != nil {
		if err == model.ErrNotFound {
			return nil, zeroerr.ErrInvalidAccount
		}
		return nil, err
	}

	// 3. Check account status
	if userInfo.Status != 1 {
		return nil, zeroerr.ErrAccountDisabled
	}

	// 4. Verify password
	if cryptx.HashPassword(in.Password, l.svcCtx.Config.Salt) != userInfo.Password {
		return nil, zeroerr.ErrInvalidPassword
	}

	// Generate tokens
	now := time.Now().Unix()
	accessToken, err := jwtx.GetToken(
		l.svcCtx.Config.JwtAuth.AccessSecret,
		now,
		l.svcCtx.Config.JwtAuth.AccessExpire,
		int64(userInfo.Id),
	)
	if err != nil {
		return nil, err
	}

	refreshToken, err := jwtx.GetToken(
		l.svcCtx.Config.JwtAuth.RefreshSecret,
		now,
		l.svcCtx.Config.JwtAuth.RefreshExpire,
		int64(userInfo.Id),
	)
	if err != nil {
		return nil, err
	}

	// Store refresh token in Redis
	err = l.svcCtx.RefreshRedis.Setex(
		fmt.Sprintf("%s%d", l.svcCtx.Config.JwtAuth.RefreshRedis.KeyPrefix, userInfo.Id),
		refreshToken,
		int(l.svcCtx.Config.JwtAuth.RefreshExpire),
	)
	if err != nil {
		return nil, err
	}

	// Record login history (async)
	go l.recordLoginHistory(userInfo.Id)

	return &user.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    l.svcCtx.Config.JwtAuth.AccessExpire,
	}, nil
}

func (l *LoginLogic) recordLoginHistory(userId uint64) {
	// Create context with timeout for async operation
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Create login record
	record := &model.LoginRecords{
		UserId:  userId,
		LoginIp: getClientIP(),
		LoginLocation: sql.NullString{
			String: getLocation(),
			Valid:  true,
		},
		DeviceType: sql.NullString{
			String: getDeviceType(),
			Valid:  true,
		},
		CreatedAt: time.Now(),
	}

	_, err := l.svcCtx.LoginRecordsModel.Insert(ctx, record)
	if err != nil {
		l.Logger.Errorf("record login history error: userId: %d, err: %v", userId, err)
	}
}

func getClientIP() string {
	return ""
}

func getLocation() string {
	return ""
}

func getDeviceType() string {
	return ""
}
