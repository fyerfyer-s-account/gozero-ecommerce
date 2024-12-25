// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.19.4
// source: marketing.proto

package marketing

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
	Marketing_CreateCoupon_FullMethodName       = "/marketing.Marketing/CreateCoupon"
	Marketing_GetCoupon_FullMethodName          = "/marketing.Marketing/GetCoupon"
	Marketing_ListCoupons_FullMethodName        = "/marketing.Marketing/ListCoupons"
	Marketing_UserCoupons_FullMethodName        = "/marketing.Marketing/UserCoupons"
	Marketing_ReceiveCoupon_FullMethodName      = "/marketing.Marketing/ReceiveCoupon"
	Marketing_VerifyCoupon_FullMethodName       = "/marketing.Marketing/VerifyCoupon"
	Marketing_UseCoupon_FullMethodName          = "/marketing.Marketing/UseCoupon"
	Marketing_CreatePromotion_FullMethodName    = "/marketing.Marketing/CreatePromotion"
	Marketing_GetPromotion_FullMethodName       = "/marketing.Marketing/GetPromotion"
	Marketing_ListPromotions_FullMethodName     = "/marketing.Marketing/ListPromotions"
	Marketing_CalculatePromotion_FullMethodName = "/marketing.Marketing/CalculatePromotion"
	Marketing_GetUserPoints_FullMethodName      = "/marketing.Marketing/GetUserPoints"
	Marketing_AddPoints_FullMethodName          = "/marketing.Marketing/AddPoints"
	Marketing_UsePoints_FullMethodName          = "/marketing.Marketing/UsePoints"
	Marketing_PointsHistory_FullMethodName      = "/marketing.Marketing/PointsHistory"
)

// MarketingClient is the client API for Marketing service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MarketingClient interface {
	// 优惠券管理
	CreateCoupon(ctx context.Context, in *CreateCouponRequest, opts ...grpc.CallOption) (*CreateCouponResponse, error)
	GetCoupon(ctx context.Context, in *GetCouponRequest, opts ...grpc.CallOption) (*GetCouponResponse, error)
	ListCoupons(ctx context.Context, in *ListCouponsRequest, opts ...grpc.CallOption) (*ListCouponsResponse, error)
	UserCoupons(ctx context.Context, in *UserCouponsRequest, opts ...grpc.CallOption) (*UserCouponsResponse, error)
	ReceiveCoupon(ctx context.Context, in *ReceiveCouponRequest, opts ...grpc.CallOption) (*ReceiveCouponResponse, error)
	VerifyCoupon(ctx context.Context, in *VerifyCouponRequest, opts ...grpc.CallOption) (*VerifyCouponResponse, error)
	UseCoupon(ctx context.Context, in *UseCouponRequest, opts ...grpc.CallOption) (*UseCouponResponse, error)
	// 促销活动
	CreatePromotion(ctx context.Context, in *CreatePromotionRequest, opts ...grpc.CallOption) (*CreatePromotionResponse, error)
	GetPromotion(ctx context.Context, in *GetPromotionRequest, opts ...grpc.CallOption) (*GetPromotionResponse, error)
	ListPromotions(ctx context.Context, in *ListPromotionsRequest, opts ...grpc.CallOption) (*ListPromotionsResponse, error)
	CalculatePromotion(ctx context.Context, in *CalculatePromotionRequest, opts ...grpc.CallOption) (*CalculatePromotionResponse, error)
	// 积分系统
	GetUserPoints(ctx context.Context, in *GetUserPointsRequest, opts ...grpc.CallOption) (*GetUserPointsResponse, error)
	AddPoints(ctx context.Context, in *AddPointsRequest, opts ...grpc.CallOption) (*AddPointsResponse, error)
	UsePoints(ctx context.Context, in *UsePointsRequest, opts ...grpc.CallOption) (*UsePointsResponse, error)
	PointsHistory(ctx context.Context, in *PointsHistoryRequest, opts ...grpc.CallOption) (*PointsHistoryResponse, error)
}

type marketingClient struct {
	cc grpc.ClientConnInterface
}

func NewMarketingClient(cc grpc.ClientConnInterface) MarketingClient {
	return &marketingClient{cc}
}

func (c *marketingClient) CreateCoupon(ctx context.Context, in *CreateCouponRequest, opts ...grpc.CallOption) (*CreateCouponResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreateCouponResponse)
	err := c.cc.Invoke(ctx, Marketing_CreateCoupon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) GetCoupon(ctx context.Context, in *GetCouponRequest, opts ...grpc.CallOption) (*GetCouponResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetCouponResponse)
	err := c.cc.Invoke(ctx, Marketing_GetCoupon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) ListCoupons(ctx context.Context, in *ListCouponsRequest, opts ...grpc.CallOption) (*ListCouponsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListCouponsResponse)
	err := c.cc.Invoke(ctx, Marketing_ListCoupons_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) UserCoupons(ctx context.Context, in *UserCouponsRequest, opts ...grpc.CallOption) (*UserCouponsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UserCouponsResponse)
	err := c.cc.Invoke(ctx, Marketing_UserCoupons_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) ReceiveCoupon(ctx context.Context, in *ReceiveCouponRequest, opts ...grpc.CallOption) (*ReceiveCouponResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ReceiveCouponResponse)
	err := c.cc.Invoke(ctx, Marketing_ReceiveCoupon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) VerifyCoupon(ctx context.Context, in *VerifyCouponRequest, opts ...grpc.CallOption) (*VerifyCouponResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(VerifyCouponResponse)
	err := c.cc.Invoke(ctx, Marketing_VerifyCoupon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) UseCoupon(ctx context.Context, in *UseCouponRequest, opts ...grpc.CallOption) (*UseCouponResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UseCouponResponse)
	err := c.cc.Invoke(ctx, Marketing_UseCoupon_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) CreatePromotion(ctx context.Context, in *CreatePromotionRequest, opts ...grpc.CallOption) (*CreatePromotionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CreatePromotionResponse)
	err := c.cc.Invoke(ctx, Marketing_CreatePromotion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) GetPromotion(ctx context.Context, in *GetPromotionRequest, opts ...grpc.CallOption) (*GetPromotionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPromotionResponse)
	err := c.cc.Invoke(ctx, Marketing_GetPromotion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) ListPromotions(ctx context.Context, in *ListPromotionsRequest, opts ...grpc.CallOption) (*ListPromotionsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(ListPromotionsResponse)
	err := c.cc.Invoke(ctx, Marketing_ListPromotions_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) CalculatePromotion(ctx context.Context, in *CalculatePromotionRequest, opts ...grpc.CallOption) (*CalculatePromotionResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(CalculatePromotionResponse)
	err := c.cc.Invoke(ctx, Marketing_CalculatePromotion_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) GetUserPoints(ctx context.Context, in *GetUserPointsRequest, opts ...grpc.CallOption) (*GetUserPointsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetUserPointsResponse)
	err := c.cc.Invoke(ctx, Marketing_GetUserPoints_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) AddPoints(ctx context.Context, in *AddPointsRequest, opts ...grpc.CallOption) (*AddPointsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(AddPointsResponse)
	err := c.cc.Invoke(ctx, Marketing_AddPoints_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) UsePoints(ctx context.Context, in *UsePointsRequest, opts ...grpc.CallOption) (*UsePointsResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UsePointsResponse)
	err := c.cc.Invoke(ctx, Marketing_UsePoints_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *marketingClient) PointsHistory(ctx context.Context, in *PointsHistoryRequest, opts ...grpc.CallOption) (*PointsHistoryResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PointsHistoryResponse)
	err := c.cc.Invoke(ctx, Marketing_PointsHistory_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MarketingServer is the server API for Marketing service.
// All implementations must embed UnimplementedMarketingServer
// for forward compatibility.
type MarketingServer interface {
	// 优惠券管理
	CreateCoupon(context.Context, *CreateCouponRequest) (*CreateCouponResponse, error)
	GetCoupon(context.Context, *GetCouponRequest) (*GetCouponResponse, error)
	ListCoupons(context.Context, *ListCouponsRequest) (*ListCouponsResponse, error)
	UserCoupons(context.Context, *UserCouponsRequest) (*UserCouponsResponse, error)
	ReceiveCoupon(context.Context, *ReceiveCouponRequest) (*ReceiveCouponResponse, error)
	VerifyCoupon(context.Context, *VerifyCouponRequest) (*VerifyCouponResponse, error)
	UseCoupon(context.Context, *UseCouponRequest) (*UseCouponResponse, error)
	// 促销活动
	CreatePromotion(context.Context, *CreatePromotionRequest) (*CreatePromotionResponse, error)
	GetPromotion(context.Context, *GetPromotionRequest) (*GetPromotionResponse, error)
	ListPromotions(context.Context, *ListPromotionsRequest) (*ListPromotionsResponse, error)
	CalculatePromotion(context.Context, *CalculatePromotionRequest) (*CalculatePromotionResponse, error)
	// 积分系统
	GetUserPoints(context.Context, *GetUserPointsRequest) (*GetUserPointsResponse, error)
	AddPoints(context.Context, *AddPointsRequest) (*AddPointsResponse, error)
	UsePoints(context.Context, *UsePointsRequest) (*UsePointsResponse, error)
	PointsHistory(context.Context, *PointsHistoryRequest) (*PointsHistoryResponse, error)
	mustEmbedUnimplementedMarketingServer()
}

// UnimplementedMarketingServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedMarketingServer struct{}

func (UnimplementedMarketingServer) CreateCoupon(context.Context, *CreateCouponRequest) (*CreateCouponResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateCoupon not implemented")
}
func (UnimplementedMarketingServer) GetCoupon(context.Context, *GetCouponRequest) (*GetCouponResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetCoupon not implemented")
}
func (UnimplementedMarketingServer) ListCoupons(context.Context, *ListCouponsRequest) (*ListCouponsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListCoupons not implemented")
}
func (UnimplementedMarketingServer) UserCoupons(context.Context, *UserCouponsRequest) (*UserCouponsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserCoupons not implemented")
}
func (UnimplementedMarketingServer) ReceiveCoupon(context.Context, *ReceiveCouponRequest) (*ReceiveCouponResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ReceiveCoupon not implemented")
}
func (UnimplementedMarketingServer) VerifyCoupon(context.Context, *VerifyCouponRequest) (*VerifyCouponResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method VerifyCoupon not implemented")
}
func (UnimplementedMarketingServer) UseCoupon(context.Context, *UseCouponRequest) (*UseCouponResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UseCoupon not implemented")
}
func (UnimplementedMarketingServer) CreatePromotion(context.Context, *CreatePromotionRequest) (*CreatePromotionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePromotion not implemented")
}
func (UnimplementedMarketingServer) GetPromotion(context.Context, *GetPromotionRequest) (*GetPromotionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPromotion not implemented")
}
func (UnimplementedMarketingServer) ListPromotions(context.Context, *ListPromotionsRequest) (*ListPromotionsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ListPromotions not implemented")
}
func (UnimplementedMarketingServer) CalculatePromotion(context.Context, *CalculatePromotionRequest) (*CalculatePromotionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CalculatePromotion not implemented")
}
func (UnimplementedMarketingServer) GetUserPoints(context.Context, *GetUserPointsRequest) (*GetUserPointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserPoints not implemented")
}
func (UnimplementedMarketingServer) AddPoints(context.Context, *AddPointsRequest) (*AddPointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddPoints not implemented")
}
func (UnimplementedMarketingServer) UsePoints(context.Context, *UsePointsRequest) (*UsePointsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UsePoints not implemented")
}
func (UnimplementedMarketingServer) PointsHistory(context.Context, *PointsHistoryRequest) (*PointsHistoryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PointsHistory not implemented")
}
func (UnimplementedMarketingServer) mustEmbedUnimplementedMarketingServer() {}
func (UnimplementedMarketingServer) testEmbeddedByValue()                   {}

// UnsafeMarketingServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MarketingServer will
// result in compilation errors.
type UnsafeMarketingServer interface {
	mustEmbedUnimplementedMarketingServer()
}

func RegisterMarketingServer(s grpc.ServiceRegistrar, srv MarketingServer) {
	// If the following call pancis, it indicates UnimplementedMarketingServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Marketing_ServiceDesc, srv)
}

func _Marketing_CreateCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateCouponRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).CreateCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_CreateCoupon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).CreateCoupon(ctx, req.(*CreateCouponRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_GetCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetCouponRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).GetCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_GetCoupon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).GetCoupon(ctx, req.(*GetCouponRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_ListCoupons_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListCouponsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).ListCoupons(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_ListCoupons_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).ListCoupons(ctx, req.(*ListCouponsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_UserCoupons_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserCouponsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).UserCoupons(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_UserCoupons_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).UserCoupons(ctx, req.(*UserCouponsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_ReceiveCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReceiveCouponRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).ReceiveCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_ReceiveCoupon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).ReceiveCoupon(ctx, req.(*ReceiveCouponRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_VerifyCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(VerifyCouponRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).VerifyCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_VerifyCoupon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).VerifyCoupon(ctx, req.(*VerifyCouponRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_UseCoupon_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UseCouponRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).UseCoupon(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_UseCoupon_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).UseCoupon(ctx, req.(*UseCouponRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_CreatePromotion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreatePromotionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).CreatePromotion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_CreatePromotion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).CreatePromotion(ctx, req.(*CreatePromotionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_GetPromotion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPromotionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).GetPromotion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_GetPromotion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).GetPromotion(ctx, req.(*GetPromotionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_ListPromotions_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListPromotionsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).ListPromotions(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_ListPromotions_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).ListPromotions(ctx, req.(*ListPromotionsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_CalculatePromotion_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CalculatePromotionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).CalculatePromotion(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_CalculatePromotion_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).CalculatePromotion(ctx, req.(*CalculatePromotionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_GetUserPoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserPointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).GetUserPoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_GetUserPoints_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).GetUserPoints(ctx, req.(*GetUserPointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_AddPoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddPointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).AddPoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_AddPoints_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).AddPoints(ctx, req.(*AddPointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_UsePoints_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UsePointsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).UsePoints(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_UsePoints_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).UsePoints(ctx, req.(*UsePointsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Marketing_PointsHistory_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PointsHistoryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MarketingServer).PointsHistory(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Marketing_PointsHistory_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MarketingServer).PointsHistory(ctx, req.(*PointsHistoryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Marketing_ServiceDesc is the grpc.ServiceDesc for Marketing service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Marketing_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "marketing.Marketing",
	HandlerType: (*MarketingServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateCoupon",
			Handler:    _Marketing_CreateCoupon_Handler,
		},
		{
			MethodName: "GetCoupon",
			Handler:    _Marketing_GetCoupon_Handler,
		},
		{
			MethodName: "ListCoupons",
			Handler:    _Marketing_ListCoupons_Handler,
		},
		{
			MethodName: "UserCoupons",
			Handler:    _Marketing_UserCoupons_Handler,
		},
		{
			MethodName: "ReceiveCoupon",
			Handler:    _Marketing_ReceiveCoupon_Handler,
		},
		{
			MethodName: "VerifyCoupon",
			Handler:    _Marketing_VerifyCoupon_Handler,
		},
		{
			MethodName: "UseCoupon",
			Handler:    _Marketing_UseCoupon_Handler,
		},
		{
			MethodName: "CreatePromotion",
			Handler:    _Marketing_CreatePromotion_Handler,
		},
		{
			MethodName: "GetPromotion",
			Handler:    _Marketing_GetPromotion_Handler,
		},
		{
			MethodName: "ListPromotions",
			Handler:    _Marketing_ListPromotions_Handler,
		},
		{
			MethodName: "CalculatePromotion",
			Handler:    _Marketing_CalculatePromotion_Handler,
		},
		{
			MethodName: "GetUserPoints",
			Handler:    _Marketing_GetUserPoints_Handler,
		},
		{
			MethodName: "AddPoints",
			Handler:    _Marketing_AddPoints_Handler,
		},
		{
			MethodName: "UsePoints",
			Handler:    _Marketing_UsePoints_Handler,
		},
		{
			MethodName: "PointsHistory",
			Handler:    _Marketing_PointsHistory_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "marketing.proto",
}
