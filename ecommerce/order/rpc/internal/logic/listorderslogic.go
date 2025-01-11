package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrdersLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrdersLogic {
	return &ListOrdersLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListOrdersLogic) ListOrders(in *order.ListOrdersRequest) (*order.ListOrdersResponse, error) {
	if in.UserId <= 0 || in.Page <= 0 || in.PageSize <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	orders, err := l.svcCtx.OrdersModel.FindByUserIdWithPage(l.ctx, uint64(in.UserId), int64(in.Status), int(in.Page), int(in.PageSize))
	if err != nil {
		return nil, err
	}

	total, err := l.svcCtx.OrdersModel.CountByUserIdAndStatus(l.ctx, uint64(in.UserId), int64(in.Status))
	if err != nil {
		return nil, err
	}

	orderList := make([]*order.Order, 0, len(orders))
	for _, o := range orders {
		// Get order items
		items, err := l.svcCtx.OrderItemsModel.FindByOrderId(l.ctx, o.Id)
		if err != nil {
			return nil, err
		}

		orderItems := make([]*order.OrderItem, 0, len(items))
		for _, item := range items {
			orderItems = append(orderItems, &order.OrderItem{
				Id:          int64(item.Id),
				OrderId:     int64(item.OrderId),
				ProductId:   int64(item.ProductId),
				SkuId:       int64(item.SkuId),
				ProductName: item.ProductName,
				SkuName:     item.SkuName,
				Price:       item.Price,
				Quantity:    item.Quantity,
				TotalAmount: item.TotalAmount,
			})
		}

		orderList = append(orderList, &order.Order{
			Id:            int64(o.Id),
			UserId:        int64(o.UserId),
			OrderNo:       o.OrderNo,
			TotalAmount:   o.TotalAmount,
			PayAmount:     o.PayAmount,
			FreightAmount: o.FreightAmount,
			Status:        int32(o.Status),
			Address:       o.Address,
			Receiver:      o.Receiver,
			Phone:         o.Phone,
			Items:         orderItems,
			CreatedAt:     o.CreatedAt.Unix(),
			UpdatedAt:     o.UpdatedAt.Unix(),
		})
	}

	return &order.ListOrdersResponse{
		Orders: orderList,
		Total:  total,
	}, nil
}
