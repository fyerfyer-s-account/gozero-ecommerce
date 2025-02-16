// Code generated by goctl. DO NOT EDIT.
// goctl 1.7.4
// Source: cart.proto

package server

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/logic"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
)

type CartServer struct {
	svcCtx *svc.ServiceContext
	cart.UnimplementedCartServer
}

func NewCartServer(svcCtx *svc.ServiceContext) *CartServer {
	return &CartServer{
		svcCtx: svcCtx,
	}
}

// 购物车操作
func (s *CartServer) AddItem(ctx context.Context, in *cart.AddItemRequest) (*cart.AddItemResponse, error) {
	l := logic.NewAddItemLogic(ctx, s.svcCtx)
	return l.AddItem(in)
}

func (s *CartServer) UpdateItem(ctx context.Context, in *cart.UpdateItemRequest) (*cart.UpdateItemResponse, error) {
	l := logic.NewUpdateItemLogic(ctx, s.svcCtx)
	return l.UpdateItem(in)
}

func (s *CartServer) RemoveItem(ctx context.Context, in *cart.RemoveItemRequest) (*cart.RemoveItemResponse, error) {
	l := logic.NewRemoveItemLogic(ctx, s.svcCtx)
	return l.RemoveItem(in)
}

func (s *CartServer) GetCart(ctx context.Context, in *cart.GetCartRequest) (*cart.GetCartResponse, error) {
	l := logic.NewGetCartLogic(ctx, s.svcCtx)
	return l.GetCart(in)
}

func (s *CartServer) ClearCart(ctx context.Context, in *cart.ClearCartRequest) (*cart.ClearCartResponse, error) {
	l := logic.NewClearCartLogic(ctx, s.svcCtx)
	return l.ClearCart(in)
}

// 商品选择
func (s *CartServer) SelectItem(ctx context.Context, in *cart.SelectItemRequest) (*cart.SelectItemResponse, error) {
	l := logic.NewSelectItemLogic(ctx, s.svcCtx)
	return l.SelectItem(in)
}

func (s *CartServer) UnselectItem(ctx context.Context, in *cart.UnselectItemRequest) (*cart.UnselectItemResponse, error) {
	l := logic.NewUnselectItemLogic(ctx, s.svcCtx)
	return l.UnselectItem(in)
}

func (s *CartServer) SelectAll(ctx context.Context, in *cart.SelectAllRequest) (*cart.SelectAllResponse, error) {
	l := logic.NewSelectAllLogic(ctx, s.svcCtx)
	return l.SelectAll(in)
}

func (s *CartServer) UnselectAll(ctx context.Context, in *cart.UnselectAllRequest) (*cart.UnselectAllResponse, error) {
	l := logic.NewUnselectAllLogic(ctx, s.svcCtx)
	return l.UnselectAll(in)
}

// 结算相关
func (s *CartServer) GetSelectedItems(ctx context.Context, in *cart.GetSelectedItemsRequest) (*cart.GetSelectedItemsResponse, error) {
	l := logic.NewGetSelectedItemsLogic(ctx, s.svcCtx)
	return l.GetSelectedItems(in)
}

func (s *CartServer) CheckStock(ctx context.Context, in *cart.CheckStockRequest) (*cart.CheckStockResponse, error) {
	l := logic.NewCheckStockLogic(ctx, s.svcCtx)
	return l.CheckStock(in)
}
