package logic

import (
	"context"
	"database/sql"
	"fmt"
	"log"
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
	log.Println("Login called")
	log.Printf("Received login request for username: %s", in.Username)

	// 1. Validate input
	if len(in.Username) == 0 || len(in.Password) == 0 {
		log.Println("Validation failed: Invalid username or password length")
		return nil, zeroerr.ErrInvalidUsername
	}

	// 2. Find user
	log.Printf("Looking up user by username or email: %s", in.Username)
	userInfo, err := l.svcCtx.UsersModel.FindOneByPhoneOrEmail(l.ctx, in.Username)
	if err != nil {
		if err == model.ErrNotFound {
			log.Println("User not found")
			return nil, zeroerr.ErrInvalidAccount
		}
		log.Printf("Error finding user: %v", err)
		return nil, err
	}

	// 3. Check account status
	log.Printf("Checking account status for user ID: %d", userInfo.Id)
	if userInfo.Status != 1 {
		log.Printf("Account disabled for user ID: %d", userInfo.Id)
		return nil, zeroerr.ErrAccountDisabled
	}

	// 4. Verify password
	log.Println("Verifying password")
	if cryptx.HashPassword(in.Password, l.svcCtx.Config.Salt) != userInfo.Password {
		log.Println("Invalid password")
		return nil, zeroerr.ErrInvalidPassword
	}

	// Check if user is admin
	role := jwtx.RoleUser
	if userInfo.IsAdmin == 1 {
		role = jwtx.RoleAdmin
	}

	// Generate tokens
	now := time.Now().Unix()
	log.Println("Generating access token")
	accessToken, err := jwtx.GetToken(
		l.svcCtx.Config.JwtAuth.AccessSecret,
		now,
		l.svcCtx.Config.JwtAuth.AccessExpire,
		int64(userInfo.Id),
		role,
	)
	if err != nil {
		log.Printf("Error generating access token: %v", err)
		return nil, zeroerr.ErrGenerateTokenFailed
	}

	log.Println("Generating refresh token")
	refreshToken, err := jwtx.GetToken(
		l.svcCtx.Config.JwtAuth.RefreshSecret,
		now,
		l.svcCtx.Config.JwtAuth.RefreshExpire,
		int64(userInfo.Id),
		role,
	)
	if err != nil {
		log.Printf("Error generating refresh token: %v", err)
		return nil, zeroerr.ErrGenerateTokenFailed
	}

	// Store refresh token in Redis
	log.Printf("Storing refresh token in Redis for user ID: %d", userInfo.Id)
	err = l.svcCtx.RefreshRedis.Setex(
		fmt.Sprintf("%s%d", l.svcCtx.Config.JwtAuth.RefreshRedis.KeyPrefix, userInfo.Id),
		refreshToken,
		int(l.svcCtx.Config.JwtAuth.RefreshExpire),
	)
	if err != nil {
		log.Printf("Error storing refresh token in Redis: %v", err)
		return nil, zeroerr.ErrStoreTokenFailed
	}

	// Record login history (async)
	go l.recordLoginHistory(userInfo.Id)

	log.Printf("Login successful for user ID: %d", userInfo.Id)
	return &user.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    l.svcCtx.Config.JwtAuth.AccessExpire,
	}, nil
}

func (l *LoginLogic) recordLoginHistory(userId uint64) {
	// Create context with timeout for async operation
	log.Printf("Start recording login history for user ID: %d", userId)
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

	log.Printf("Inserting login history record for user ID: %d", userId)
	_, err := l.svcCtx.LoginRecordsModel.Insert(ctx, record)
	if err != nil {
		log.Printf("Error recording login history for user ID: %d, error: %v", userId, err)
	}
}

func getClientIP() string {
	log.Println("Fetching client IP")
	// Placeholder for actual implementation
	return "127.0.0.1"
}

func getLocation() string {
	log.Println("Fetching login location")
	// Placeholder for actual implementation
	return "Unknown Location"
}

func getDeviceType() string {
	log.Println("Fetching device type")
	// Placeholder for actual implementation
	return "Unknown Device"
}
