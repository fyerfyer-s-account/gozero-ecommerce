package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetOrderLogic) GetOrder(in *order.GetOrderRequest) (*order.GetOrderResponse, error) {
	if len(in.OrderNo) == 0 {
		return nil, zeroerr.ErrOrderNoEmpty
	}

	orderInfo, err := l.svcCtx.OrdersModel.FindOneByOrderNo(l.ctx, in.OrderNo)
	if err != nil {
		return nil, err
	}

	items, err := l.svcCtx.OrderItemsModel.FindByOrderId(l.ctx, orderInfo.Id)
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

	// Get payment info
	payment, err := l.svcCtx.OrderPaymentsModel.FindByOrderId(l.ctx, orderInfo.Id)
	var paymentInfo *order.PaymentInfo
	if err == nil && payment != nil {
		paymentInfo = &order.PaymentInfo{
			Id:            int64(payment.Id),
			OrderId:       int64(payment.OrderId),
			PaymentNo:     payment.PaymentNo,
			PaymentMethod: int32(payment.PaymentMethod),
			Amount:        payment.Amount,
			Status:        int32(payment.Status),
			PayTime:       payment.PayTime.Time.Unix(),
		}
	}

	// Get shipping info
	shipping, err := l.svcCtx.OrderShippingModel.FindByOrderId(l.ctx, orderInfo.Id)
	var shippingInfo *order.ShippingInfo
	if err == nil && shipping != nil {
		shippingInfo = &order.ShippingInfo{
			Id:          int64(shipping.Id),
			OrderId:     int64(shipping.OrderId),
			ShippingNo:  shipping.ShippingNo.String,
			Company:     shipping.Company.String,
			Status:      int32(shipping.Status),
			ShipTime:    shipping.ShipTime.Time.Unix(),
			ReceiveTime: shipping.ReceiveTime.Time.Unix(),
		}
	}

	return &order.GetOrderResponse{
		Order: &order.Order{
			Id:            int64(orderInfo.Id),
			UserId:        int64(orderInfo.UserId),
			OrderNo:       orderInfo.OrderNo,
			TotalAmount:   orderInfo.TotalAmount,
			PayAmount:     orderInfo.PayAmount,
			FreightAmount: orderInfo.FreightAmount,
			Status:        int32(orderInfo.Status),
			Address:       orderInfo.Address,
			Receiver:      orderInfo.Receiver,
			Phone:         orderInfo.Phone,
			Items:         orderItems,
			Payment:       paymentInfo,
			Shipping:      shippingInfo,
			CreatedAt:     orderInfo.CreatedAt.Unix(),
			UpdatedAt:     orderInfo.UpdatedAt.Unix(),
		},
	}, nil
}
