package logic

import (
	"context"
	"log"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 获取用户信息
func (l *GetUserInfoLogic) GetUserInfo(in *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	log.Println("GetUserInfo called")
	log.Printf("Received request for user ID: %d", in.UserId)

	// 1. Get user basic info
	log.Println("Fetching user basic information")
	userInfo, err := l.svcCtx.UsersModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil {
		log.Printf("Error fetching user info: %v", err)
		return nil, zeroerr.ErrUserNotFound
	}
	log.Printf("User info found for user ID: %d", in.UserId)

	// 2. Get wallet info
	log.Println("Fetching wallet information")
	wallet, err := l.svcCtx.WalletAccountsModel.FindOne(l.ctx, uint64(in.UserId))
	if err != nil && err != model.ErrNotFound {
		log.Printf("Error fetching wallet info: %v", err)
		return nil, err
	}
	if err == model.ErrNotFound {
		log.Printf("No wallet found for user ID: %d", in.UserId)
	} else {
		log.Printf("Wallet info found for user ID: %d, balance: %.2f", in.UserId, wallet.Balance)
	}

	// Prepare response
	log.Println("Preparing response data")
	response := &user.GetUserInfoResponse{
		UserInfo: &user.UserInfo{
			UserId:        int64(userInfo.Id),
			Username:      userInfo.Username,
			Nickname:      userInfo.Nickname.String,
			Avatar:        userInfo.Avatar.String,
			Phone:         userInfo.Phone.String,
			Email:         userInfo.Email.String,
			Gender:        userInfo.Gender,
			WalletBalance: wallet.Balance,
			CreatedAt:     userInfo.CreatedAt.Unix(),
			UpdatedAt:     userInfo.UpdatedAt.Unix(),
		},
	}

	log.Println("Returning user info response")
	return response, nil
}
