package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateOrderLogic {
	return &CreateOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 订单管理
func (l *CreateOrderLogic) CreateOrder(in *order.CreateOrderRequest) (*order.CreateOrderResponse, error) {
	if in.UserId <= 0 || len(in.Items) == 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	orderNo := fmt.Sprintf("%d%d", time.Now().UnixNano(), in.UserId)

	// Create order record
	orderData := &model.Orders{
		OrderNo:   orderNo,
		UserId:    uint64(in.UserId),
		Status:    1, // Pending payment
		Address:   in.Address,
		Receiver:  in.Receiver,
		Phone:     in.Phone,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Calculate total amount and create order items
	var totalAmount float64
	orderItems := make([]*model.OrderItems, 0, len(in.Items))
	for _, item := range in.Items {
		// Here you should call ProductRpc to get product info
		// For now we'll just use placeholder values
		orderItem := &model.OrderItems{
			OrderId:   0, // Will be set after order creation
			ProductId: uint64(item.ProductId),
			SkuId:     uint64(item.SkuId),
			Quantity:  int64(item.Quantity),
			CreatedAt: time.Now(),
		}
		orderItems = append(orderItems, orderItem)
		totalAmount += float64(item.Quantity) * 100 // Placeholder price
	}

	orderData.TotalAmount = totalAmount
	orderData.PayAmount = totalAmount

	// Create order
	orderId, err := l.svcCtx.OrdersModel.CreateOrder(l.ctx, orderData)
	if err != nil {
		return nil, zeroerr.ErrOrderCreateFailed
	}

	// Update order items with order ID
	for _, item := range orderItems {
		item.OrderId = orderId
	}

	// Batch insert order items
	err = l.svcCtx.OrderItemsModel.BatchInsert(l.ctx, orderItems)
	if err != nil {
		return nil, zeroerr.ErrOrderCreateFailed
	}

	// Create shipping record
	shipping := &model.OrderShipping{
		OrderId:   orderId,
		Status:    0, // Pending
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	_, err = l.svcCtx.OrderShippingModel.Insert(l.ctx, shipping)
	if err != nil {
		return nil, zeroerr.ErrOrderCreateFailed
	}

	return &order.CreateOrderResponse{
		OrderNo:   orderNo,
		PayAmount: totalAmount,
	}, nil
}
