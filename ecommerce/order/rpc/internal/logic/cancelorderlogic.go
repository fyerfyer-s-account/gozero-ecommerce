package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type CancelOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCancelOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CancelOrderLogic {
	return &CancelOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CancelOrderLogic) CancelOrder(in *order.CancelOrderRequest) (*order.CancelOrderResponse, error) {
	if len(in.OrderNo) == 0 {
		return nil, zeroerr.ErrOrderInvalidParam
	}

	orderInfo, err := l.svcCtx.OrdersModel.FindOneByOrderNo(l.ctx, in.OrderNo)
	if err != nil {
		return nil, err
	}

	if orderInfo.Status != 1 {
		return nil, zeroerr.ErrOrderStatusInvalid
	}

	err = l.svcCtx.OrdersModel.UpdateStatus(l.ctx, orderInfo.Id, 5) // 5: Canceled
	if err != nil {
		return nil, err
	}

	// Get payment & shipping number

	payment, _ := l.svcCtx.OrderPaymentsModel.FindByOrderId(l.ctx, orderInfo.Id)

	shipping, _ := l.svcCtx.OrderShippingModel.FindByOrderId(l.ctx, orderInfo.Id)

	// Use status change event since we don't have dedicated cancel publisher
	statusEvent := &types.OrderStatusChangedEvent{
		OldStatus:  1, // From pending payment
		NewStatus:  5, // To canceled
		EventType:  types.OrderStatusCanceled,
		PaymentNo:  payment.PaymentNo,
		ShippingNo: shipping.ShippingNo.String,
		Reason:     in.Reason,
	}

	if err := l.svcCtx.Producer.PublishStatusChanged(l.ctx, statusEvent); err != nil {
		logx.Errorf("failed to publish order cancel status event: %v", err)
	}

	return &order.CancelOrderResponse{
		Success: true,
	}, nil
}

func convertToEventItems(items []*model.OrderItems) []types.OrderItem {
	result := make([]types.OrderItem, len(items))
	for i, item := range items {
		result[i] = types.OrderItem{
			SkuID:     int64(item.SkuId),
			ProductID: int64(item.ProductId),
			Quantity:  int32(item.Quantity),
			Price:     item.Price,
		}
	}
	return result
}
