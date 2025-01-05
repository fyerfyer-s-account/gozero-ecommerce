// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.4
// source: user.proto

package user

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	User_Register_FullMethodName         = "/user.User/Register"
	User_Login_FullMethodName            = "/user.User/Login"
	User_Logout_FullMethodName           = "/user.User/Logout"
	User_GetUserInfo_FullMethodName      = "/user.User/GetUserInfo"
	User_GetUserAddresses_FullMethodName = "/user.User/GetUserAddresses"
	User_GetTransactions_FullMethodName  = "/user.User/GetTransactions"
	User_UpdateUserInfo_FullMethodName   = "/user.User/UpdateUserInfo"
	User_ChangePassword_FullMethodName   = "/user.User/ChangePassword"
	User_ResetPassword_FullMethodName    = "/user.User/ResetPassword"
	User_AddAddress_FullMethodName       = "/user.User/AddAddress"
	User_UpdateAddress_FullMethodName    = "/user.User/UpdateAddress"
	User_DeleteAddress_FullMethodName    = "/user.User/DeleteAddress"
	User_GetWallet_FullMethodName        = "/user.User/GetWallet"
	User_RechargeWallet_FullMethodName   = "/user.User/RechargeWallet"
	User_WithdrawWallet_FullMethodName   = "/user.User/WithdrawWallet"
)

// UserClient is the client API for User service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
//
// 用户服务接口定义
type UserClient interface {
	// 用户注册
	Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error)
	// 用户登录
	Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error)
	Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error)
	// 获取用户信息
	GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...grpc.CallOption) (*GetUserInfoResponse, error)
	GetUserAddresses(ctx context.Context, in *GetUserAddressesRequest, opts ...grpc.CallOption) (*GetUserAddressesResponse, error)
	GetTransactions(ctx context.Context, in *GetTransactionsRequest, opts ...grpc.CallOption) (*GetTransactionsResponse, error)
	// 更新用户信息
	UpdateUserInfo(ctx context.Context, in *UpdateUserInfoRequest, opts ...grpc.CallOption) (*UpdateUserInfoResponse, error)
	// 修改密码
	ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error)
	// 重置密码
	ResetPassword(ctx context.Context, in *ResetPasswordRequest, opts ...grpc.CallOption) (*ResetPasswordResponse, error)
	// 地址管理
	AddAddress(ctx context.Context, in *AddAddressRequest, opts ...grpc.CallOption) (*AddAddressResponse, error)
	UpdateAddress(ctx context.Context, in *UpdateAddressRequest, opts ...grpc.CallOption) (*UpdateAddressResponse, error)
	DeleteAddress(ctx context.Context, in *DeleteAddressRequest, opts ...grpc.CallOption) (*DeleteAddressResponse, error)
	// 钱包操作
	GetWallet(ctx context.Context, in *GetWalletRequest, opts ...grpc.CallOption) (*GetWalletResponse, error)
	RechargeWallet(ctx context.Context, in *RechargeWalletRequest, opts ...grpc.CallOption) (*RechargeWalletResponse, error)
	WithdrawWallet(ctx context.Context, in *WithdrawWalletRequest, opts ...grpc.CallOption) (*WithdrawWalletResponse, error)
}

type userClient struct {
	cc grpc.ClientConnInterface
}

func NewUserClient(cc grpc.ClientConnInterface) UserClient {
	return &userClient{cc}
}

func (c *userClient) Register(ctx context.Context, in *RegisterRequest, opts ...grpc.CallOption) (*RegisterResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RegisterResponse)
	err := c.cc.Invoke(ctx, User_Register_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) Login(ctx context.Context, in *LoginRequest, opts ...grpc.CallOption) (*LoginResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoginResponse)
	err := c.cc.Invoke(ctx, User_Login_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) Logout(ctx context.Context, in *LogoutRequest, opts ...grpc.CallOption) (*LogoutResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LogoutResponse)
	err := c.cc.Invoke(ctx, User_Logout_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetUserInfo(ctx context.Context, in *GetUserInfoRequest, opts ...grpc.CallOption) (*GetUserInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserInfoResponse)
	err := c.cc.Invoke(ctx, User_GetUserInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetUserAddresses(ctx context.Context, in *GetUserAddressesRequest, opts ...grpc.CallOption) (*GetUserAddressesResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserAddressesResponse)
	err := c.cc.Invoke(ctx, User_GetUserAddresses_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetTransactions(ctx context.Context, in *GetTransactionsRequest, opts ...grpc.CallOption) (*GetTransactionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetTransactionsResponse)
	err := c.cc.Invoke(ctx, User_GetTransactions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) UpdateUserInfo(ctx context.Context, in *UpdateUserInfoRequest, opts ...grpc.CallOption) (*UpdateUserInfoResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateUserInfoResponse)
	err := c.cc.Invoke(ctx, User_UpdateUserInfo_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) ChangePassword(ctx context.Context, in *ChangePasswordRequest, opts ...grpc.CallOption) (*ChangePasswordResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ChangePasswordResponse)
	err := c.cc.Invoke(ctx, User_ChangePassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) ResetPassword(ctx context.Context, in *ResetPasswordRequest, opts ...grpc.CallOption) (*ResetPasswordResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ResetPasswordResponse)
	err := c.cc.Invoke(ctx, User_ResetPassword_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) AddAddress(ctx context.Context, in *AddAddressRequest, opts ...grpc.CallOption) (*AddAddressResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddAddressResponse)
	err := c.cc.Invoke(ctx, User_AddAddress_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) UpdateAddress(ctx context.Context, in *UpdateAddressRequest, opts ...grpc.CallOption) (*UpdateAddressResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateAddressResponse)
	err := c.cc.Invoke(ctx, User_UpdateAddress_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) DeleteAddress(ctx context.Context, in *DeleteAddressRequest, opts ...grpc.CallOption) (*DeleteAddressResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(DeleteAddressResponse)
	err := c.cc.Invoke(ctx, User_DeleteAddress_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) GetWallet(ctx context.Context, in *GetWalletRequest, opts ...grpc.CallOption) (*GetWalletResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetWalletResponse)
	err := c.cc.Invoke(ctx, User_GetWallet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) RechargeWallet(ctx context.Context, in *RechargeWalletRequest, opts ...grpc.CallOption) (*RechargeWalletResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(RechargeWalletResponse)
	err := c.cc.Invoke(ctx, User_RechargeWallet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userClient) WithdrawWallet(ctx context.Context, in *WithdrawWalletRequest, opts ...grpc.CallOption) (*WithdrawWalletResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(WithdrawWalletResponse)
	err := c.cc.Invoke(ctx, User_WithdrawWallet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServer is the server API for User service.
// All implementations must embed UnimplementedUserServer
// for forward compatibility.
//
// 用户服务接口定义
type UserServer interface {
	// 用户注册
	Register(context.Context, *RegisterRequest) (*RegisterResponse, error)
	// 用户登录
	Login(context.Context, *LoginRequest) (*LoginResponse, error)
	Logout(context.Context, *LogoutRequest) (*LogoutResponse, error)
	// 获取用户信息
	GetUserInfo(context.Context, *GetUserInfoRequest) (*GetUserInfoResponse, error)
	GetUserAddresses(context.Context, *GetUserAddressesRequest) (*GetUserAddressesResponse, error)
	GetTransactions(context.Context, *GetTransactionsRequest) (*GetTransactionsResponse, error)
	// 更新用户信息
	UpdateUserInfo(context.Context, *UpdateUserInfoRequest) (*UpdateUserInfoResponse, error)
	// 修改密码
	ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponse, error)
	// 重置密码
	ResetPassword(context.Context, *ResetPasswordRequest) (*ResetPasswordResponse, error)
	// 地址管理
	AddAddress(context.Context, *AddAddressRequest) (*AddAddressResponse, error)
	UpdateAddress(context.Context, *UpdateAddressRequest) (*UpdateAddressResponse, error)
	DeleteAddress(context.Context, *DeleteAddressRequest) (*DeleteAddressResponse, error)
	// 钱包操作
	GetWallet(context.Context, *GetWalletRequest) (*GetWalletResponse, error)
	RechargeWallet(context.Context, *RechargeWalletRequest) (*RechargeWalletResponse, error)
	WithdrawWallet(context.Context, *WithdrawWalletRequest) (*WithdrawWalletResponse, error)
	mustEmbedUnimplementedUserServer()
}

// UnimplementedUserServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedUserServer struct{}

func (UnimplementedUserServer) Register(context.Context, *RegisterRequest) (*RegisterResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedUserServer) Login(context.Context, *LoginRequest) (*LoginResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Login not implemented")
}
func (UnimplementedUserServer) Logout(context.Context, *LogoutRequest) (*LogoutResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Logout not implemented")
}
func (UnimplementedUserServer) GetUserInfo(context.Context, *GetUserInfoRequest) (*GetUserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserInfo not implemented")
}
func (UnimplementedUserServer) GetUserAddresses(context.Context, *GetUserAddressesRequest) (*GetUserAddressesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserAddresses not implemented")
}
func (UnimplementedUserServer) GetTransactions(context.Context, *GetTransactionsRequest) (*GetTransactionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetTransactions not implemented")
}
func (UnimplementedUserServer) UpdateUserInfo(context.Context, *UpdateUserInfoRequest) (*UpdateUserInfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateUserInfo not implemented")
}
func (UnimplementedUserServer) ChangePassword(context.Context, *ChangePasswordRequest) (*ChangePasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ChangePassword not implemented")
}
func (UnimplementedUserServer) ResetPassword(context.Context, *ResetPasswordRequest) (*ResetPasswordResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ResetPassword not implemented")
}
func (UnimplementedUserServer) AddAddress(context.Context, *AddAddressRequest) (*AddAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddAddress not implemented")
}
func (UnimplementedUserServer) UpdateAddress(context.Context, *UpdateAddressRequest) (*UpdateAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateAddress not implemented")
}
func (UnimplementedUserServer) DeleteAddress(context.Context, *DeleteAddressRequest) (*DeleteAddressResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DeleteAddress not implemented")
}
func (UnimplementedUserServer) GetWallet(context.Context, *GetWalletRequest) (*GetWalletResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetWallet not implemented")
}
func (UnimplementedUserServer) RechargeWallet(context.Context, *RechargeWalletRequest) (*RechargeWalletResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RechargeWallet not implemented")
}
func (UnimplementedUserServer) WithdrawWallet(context.Context, *WithdrawWalletRequest) (*WithdrawWalletResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method WithdrawWallet not implemented")
}
func (UnimplementedUserServer) mustEmbedUnimplementedUserServer() {}
func (UnimplementedUserServer) testEmbeddedByValue()              {}

// UnsafeUserServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServer will
// result in compilation errors.
type UnsafeUserServer interface {
	mustEmbedUnimplementedUserServer()
}

func RegisterUserServer(s grpc.ServiceRegistrar, srv UserServer) {
	// If the following call pancis, it indicates UnimplementedUserServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&User_ServiceDesc, srv)
}

func _User_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RegisterRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_Register_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Register(ctx, req.(*RegisterRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_Login_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoginRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Login(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_Login_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Login(ctx, req.(*LoginRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_Logout_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LogoutRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).Logout(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_Logout_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).Logout(ctx, req.(*LogoutRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetUserInfo(ctx, req.(*GetUserInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetUserAddresses_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserAddressesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetUserAddresses(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetUserAddresses_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetUserAddresses(ctx, req.(*GetUserAddressesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetTransactions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetTransactionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetTransactions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetTransactions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetTransactions(ctx, req.(*GetTransactionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_UpdateUserInfo_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserInfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).UpdateUserInfo(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_UpdateUserInfo_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).UpdateUserInfo(ctx, req.(*UpdateUserInfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_ChangePassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ChangePasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).ChangePassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_ChangePassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).ChangePassword(ctx, req.(*ChangePasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_ResetPassword_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ResetPasswordRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).ResetPassword(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_ResetPassword_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).ResetPassword(ctx, req.(*ResetPasswordRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_AddAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).AddAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_AddAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).AddAddress(ctx, req.(*AddAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_UpdateAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).UpdateAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_UpdateAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).UpdateAddress(ctx, req.(*UpdateAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_DeleteAddress_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAddressRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).DeleteAddress(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_DeleteAddress_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).DeleteAddress(ctx, req.(*DeleteAddressRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_GetWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).GetWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_GetWallet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).GetWallet(ctx, req.(*GetWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_RechargeWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RechargeWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).RechargeWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_RechargeWallet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).RechargeWallet(ctx, req.(*RechargeWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _User_WithdrawWallet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(WithdrawWalletRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServer).WithdrawWallet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: User_WithdrawWallet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServer).WithdrawWallet(ctx, req.(*WithdrawWalletRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// User_ServiceDesc is the grpc.ServiceDesc for User service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var User_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "user.User",
	HandlerType: (*UserServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _User_Register_Handler,
		},
		{
			MethodName: "Login",
			Handler:    _User_Login_Handler,
		},
		{
			MethodName: "Logout",
			Handler:    _User_Logout_Handler,
		},
		{
			MethodName: "GetUserInfo",
			Handler:    _User_GetUserInfo_Handler,
		},
		{
			MethodName: "GetUserAddresses",
			Handler:    _User_GetUserAddresses_Handler,
		},
		{
			MethodName: "GetTransactions",
			Handler:    _User_GetTransactions_Handler,
		},
		{
			MethodName: "UpdateUserInfo",
			Handler:    _User_UpdateUserInfo_Handler,
		},
		{
			MethodName: "ChangePassword",
			Handler:    _User_ChangePassword_Handler,
		},
		{
			MethodName: "ResetPassword",
			Handler:    _User_ResetPassword_Handler,
		},
		{
			MethodName: "AddAddress",
			Handler:    _User_AddAddress_Handler,
		},
		{
			MethodName: "UpdateAddress",
			Handler:    _User_UpdateAddress_Handler,
		},
		{
			MethodName: "DeleteAddress",
			Handler:    _User_DeleteAddress_Handler,
		},
		{
			MethodName: "GetWallet",
			Handler:    _User_GetWallet_Handler,
		},
		{
			MethodName: "RechargeWallet",
			Handler:    _User_RechargeWallet_Handler,
		},
		{
			MethodName: "WithdrawWallet",
			Handler:    _User_WithdrawWallet_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "user.proto",
}
