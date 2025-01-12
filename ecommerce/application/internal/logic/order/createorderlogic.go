package order

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cartclient"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/orderservice"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/user/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateOrderLogic) CreateOrder(req *types.CreateOrderReq) (*types.Order, error) {
	userId := l.ctx.Value("userId").(int64)

	// Get selected cart items first
	cartResp, err := l.svcCtx.CartRpc.GetSelectedItems(l.ctx, &cartclient.GetSelectedItemsRequest{
		UserId: userId,
	})
	if err != nil {
		return nil, err
	}

	// Check if cart items are empty
	if len(cartResp.Items) == 0 {
		return nil, zeroerr.ErrInvalidParam // Ensure that we return an error for empty cart
	}

	// Get address details
	addressResp, err := l.svcCtx.UserRpc.GetAddress(l.ctx, &user.GetAddressRequest{
		AddressId: req.AddressId,
	})
	if err != nil {
		return nil, err
	}

	// Convert cart items to order items
	orderItems := make([]*order.OrderItemRequest, 0, len(cartResp.Items))
	for _, item := range cartResp.Items {
		// Verify product and SKU still exist
		_, err := l.svcCtx.ProductRpc.GetSku(l.ctx, &product.GetSkuRequest{
			Id: item.SkuId,
		})
		if err != nil {
			return nil, err
		}

		orderItems = append(orderItems, &order.OrderItemRequest{
			ProductId: item.ProductId,
			SkuId:     item.SkuId,
			Quantity:  int32(item.Quantity),
		})
	}

	// Create order
	orderResp, err := l.svcCtx.OrderRpc.CreateOrder(l.ctx, &order.CreateOrderRequest{
		UserId:   userId,
		Address:  addressResp.Address.DetailAddress,
		Receiver: addressResp.Address.ReceiverName,
		Phone:    addressResp.Address.ReceiverPhone,
		Items:    orderItems,
	})
	if err != nil {
		return nil, err
	}

	// Get created order details
	getResp, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderRequest{
		OrderNo: orderResp.OrderNo,
	})
	if err != nil {
		return nil, err
	}

	return convertOrder(getResp.Order, req.Note), nil
}


func convertOrder(o *orderservice.Order, note string) *types.Order {
	items := make([]types.OrderItem, 0, len(o.Items))
	for _, item := range o.Items {
		items = append(items, types.OrderItem{
			Id:          item.Id,
			ProductId:   item.ProductId,
			ProductName: item.ProductName,
			SkuId:       item.SkuId,
			SkuName:     item.SkuName,
			Price:       item.Price,
			Quantity:    int32(item.Quantity),
			Amount:      item.TotalAmount,
		})
	}

	order := &types.Order{
		Id:            o.Id,
		OrderNo:       o.OrderNo,
		UserId:        o.UserId,
		Status:        o.Status,
		TotalAmount:   o.TotalAmount,
		PayAmount:     o.PayAmount,
		FreightAmount: o.FreightAmount,
		Items:         items,
		CreatedAt:     o.CreatedAt,
	}

	if note != "" {
		order.Note = note
	}

	return order
}
