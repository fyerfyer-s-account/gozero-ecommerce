package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type PayOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPayOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PayOrderLogic {
	return &PayOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PayOrderLogic) PayOrder(in *order.PayOrderRequest) (*order.PayOrderResponse, error) {
    if len(in.OrderNo) == 0 || in.PaymentMethod <= 0 || in.PaymentMethod > 3 {
        return nil, zeroerr.ErrInvalidParam
    }

    orderInfo, err := l.svcCtx.OrdersModel.FindOneByOrderNo(l.ctx, in.OrderNo)
    if err != nil {
        return nil, err
    }

    if orderInfo.Status != 1 {
        return nil, zeroerr.ErrOrderStatusNotAllowed
    }

    paymentNo := fmt.Sprintf("PAY%d%d", time.Now().UnixNano(), orderInfo.UserId)
    
    payment := &model.OrderPayments{
        OrderId:       orderInfo.Id,
        PaymentNo:     paymentNo,
        PaymentMethod: int64(in.PaymentMethod),
        Amount:        orderInfo.PayAmount,
        Status:        0, // Unpaid
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    _, err = l.svcCtx.OrderPaymentsModel.Insert(l.ctx, payment)
    if err != nil {
        return nil, err
    }

    // Generate mock pay URL based on payment method
    var payURL string
    switch in.PaymentMethod {
    case 1:
        payURL = fmt.Sprintf("https://wx.pay.com/%s", paymentNo)
    case 2:
        payURL = fmt.Sprintf("https://alipay.com/%s", paymentNo)
    case 3:
        payURL = fmt.Sprintf("https://balance.pay.com/%s", paymentNo)
    }

    // After payment record creation, publish events
    // 1. Order paid event
    paidEvent := &types.OrderPaidEvent{
        OrderEvent: types.OrderEvent{
            Type:      types.OrderPaid,
            OrderNo:   orderInfo.OrderNo,
            UserID:    int64(orderInfo.UserId),
            Timestamp: time.Now(),
        },
        PaymentNo:     paymentNo,
        PaymentMethod: in.PaymentMethod,
        PayAmount:     orderInfo.PayAmount,
    }

    if err := l.svcCtx.Producer.PublishOrderPaid(l.ctx, paidEvent); err != nil {
        logx.Errorf("failed to publish order paid event: %v", err)
        // Don't return error as payment is already created
    }

    // 2. Status change event
    statusEvent := &types.OrderStatusChangedEvent{
        OrderEvent: types.OrderEvent{
            Type:      types.OrderEventType(types.OrderStatusPaid),
            OrderNo:   orderInfo.OrderNo,
            UserID:    int64(orderInfo.UserId),
            Timestamp: time.Now(),
        },
        OldStatus:  1, // Pending payment
        NewStatus:  2, // Paid, waiting for shipment
        EventType:  types.OrderStatusPaid,
        PaymentNo:  paymentNo,
    }

    if err := l.svcCtx.Producer.PublishStatusChanged(l.ctx, statusEvent); err != nil {
        logx.Errorf("failed to publish status change event: %v", err)
    }

    return &order.PayOrderResponse{
        PaymentNo: paymentNo,
        PayUrl:    payURL,
    }, nil
}
