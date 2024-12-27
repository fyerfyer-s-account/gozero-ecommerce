// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.3
// Source: user.proto

package server

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/logic"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"
)

type UserServer struct {
	svcCtx *svc.ServiceContext
	user.UnimplementedUserServer
}

func NewUserServer(svcCtx *svc.ServiceContext) *UserServer {
	return &UserServer{
		svcCtx: svcCtx,
	}
}

// 用户注册
func (s *UserServer) Register(ctx context.Context, in *user.RegisterRequest) (*user.RegisterResponse, error) {
	l := logic.NewRegisterLogic(ctx, s.svcCtx)
	return l.Register(in)
}

// 用户登录
func (s *UserServer) Login(ctx context.Context, in *user.LoginRequest) (*user.LoginResponse, error) {
	l := logic.NewLoginLogic(ctx, s.svcCtx)
	return l.Login(in)
}

// 获取用户信息
func (s *UserServer) GetUserInfo(ctx context.Context, in *user.GetUserInfoRequest) (*user.GetUserInfoResponse, error) {
	l := logic.NewGetUserInfoLogic(ctx, s.svcCtx)
	return l.GetUserInfo(in)
}

func (s *UserServer) GetUserAddresses(ctx context.Context, in *user.GetUserAddressesRequest) (*user.GetUserAddressesResponse, error) {
	l := logic.NewGetUserAddressesLogic(ctx, s.svcCtx)
	return l.GetUserAddresses(in)
}

func (s *UserServer) GetTransactions(ctx context.Context, in *user.GetTransactionsRequest) (*user.GetTransactionsResponse, error) {
	l := logic.NewGetTransactionsLogic(ctx, s.svcCtx)
	return l.GetTransactions(in)
}

// 更新用户信息
func (s *UserServer) UpdateUserInfo(ctx context.Context, in *user.UpdateUserInfoRequest) (*user.UpdateUserInfoResponse, error) {
	l := logic.NewUpdateUserInfoLogic(ctx, s.svcCtx)
	return l.UpdateUserInfo(in)
}

// 修改密码
func (s *UserServer) ChangePassword(ctx context.Context, in *user.ChangePasswordRequest) (*user.ChangePasswordResponse, error) {
	l := logic.NewChangePasswordLogic(ctx, s.svcCtx)
	return l.ChangePassword(in)
}

// 重置密码
func (s *UserServer) ResetPassword(ctx context.Context, in *user.ResetPasswordRequest) (*user.ResetPasswordResponse, error) {
	l := logic.NewResetPasswordLogic(ctx, s.svcCtx)
	return l.ResetPassword(in)
}

// 地址管理
func (s *UserServer) AddAddress(ctx context.Context, in *user.AddAddressRequest) (*user.AddAddressResponse, error) {
	l := logic.NewAddAddressLogic(ctx, s.svcCtx)
	return l.AddAddress(in)
}

func (s *UserServer) UpdateAddress(ctx context.Context, in *user.UpdateAddressRequest) (*user.UpdateAddressResponse, error) {
	l := logic.NewUpdateAddressLogic(ctx, s.svcCtx)
	return l.UpdateAddress(in)
}

func (s *UserServer) DeleteAddress(ctx context.Context, in *user.DeleteAddressRequest) (*user.DeleteAddressResponse, error) {
	l := logic.NewDeleteAddressLogic(ctx, s.svcCtx)
	return l.DeleteAddress(in)
}

// 钱包操作
func (s *UserServer) GetWallet(ctx context.Context, in *user.GetWalletRequest) (*user.GetWalletResponse, error) {
	l := logic.NewGetWalletLogic(ctx, s.svcCtx)
	return l.GetWallet(in)
}

func (s *UserServer) RechargeWallet(ctx context.Context, in *user.RechargeWalletRequest) (*user.RechargeWalletResponse, error) {
	l := logic.NewRechargeWalletLogic(ctx, s.svcCtx)
	return l.RechargeWallet(in)
}

func (s *UserServer) WithdrawWallet(ctx context.Context, in *user.WithdrawWalletRequest) (*user.WithdrawWalletResponse, error) {
	l := logic.NewWithdrawWalletLogic(ctx, s.svcCtx)
	return l.WithdrawWallet(in)
}
