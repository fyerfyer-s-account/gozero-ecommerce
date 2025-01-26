package logic

import (
	"context"
	"database/sql"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type ReceiveOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewReceiveOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ReceiveOrderLogic {
	return &ReceiveOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ReceiveOrderLogic) ReceiveOrder(in *order.ReceiveOrderRequest) (*order.ReceiveOrderResponse, error) {
    if len(in.OrderNo) == 0 {
        return nil, zeroerr.ErrOrderNoEmpty
    }

    orderInfo, err := l.svcCtx.OrdersModel.FindOneByOrderNo(l.ctx, in.OrderNo)
    if err != nil {
        return nil, err
    }

    if orderInfo.Status != 3 { // 3: shipped
        return nil, zeroerr.ErrOrderStatusNotAllowed
    }

    // Update order status to completed
    err = l.svcCtx.OrdersModel.UpdateStatus(l.ctx, orderInfo.Id, 4)
    if err != nil {
        return nil, err
    }

    // Update shipping status
    shipping, err := l.svcCtx.OrderShippingModel.FindByOrderId(l.ctx, orderInfo.Id)
    if err != nil {
        return nil, err
    }

    shipping.Status = 2 // Received
    shipping.ReceiveTime = sql.NullTime {
		Time: time.Now(),
		Valid: true,
	}
    err = l.svcCtx.OrderShippingModel.Update(l.ctx, shipping)
    if err != nil {
        return nil, err
    }

    // After status updates, publish events
    // 1. Order completed event
    completedEvent := &types.OrderCompletedEvent{
        OrderEvent: types.OrderEvent{
            Type:      types.OrderCompleted,
            OrderNo:   orderInfo.OrderNo,
            UserID:    int64(orderInfo.UserId),
            Timestamp: time.Now(),
        },
        ReceiveTime: shipping.ReceiveTime.Time,
    }

    if err := l.svcCtx.Producer.PublishOrderCompleted(l.ctx, completedEvent); err != nil {
        logx.Errorf("failed to publish order completed event: %v", err)
    }

    // 2. Status change event
    statusEvent := &types.OrderStatusChangedEvent{
        OrderEvent: types.OrderEvent{
            Type:      types.OrderEventType(types.OrderStatusReceived),
            OrderNo:   orderInfo.OrderNo,
            UserID:    int64(orderInfo.UserId),
            Timestamp: time.Now(),
        },
        OldStatus:  3, // Shipped
        NewStatus:  4, // Completed
        EventType:  types.OrderStatusReceived,
        ShippingNo: shipping.ShippingNo.String,
    }

    if err := l.svcCtx.Producer.PublishStatusChanged(l.ctx, statusEvent); err != nil {
        logx.Errorf("failed to publish status change event: %v", err)
    }

    return &order.ReceiveOrderResponse{
        Success: true,
    }, nil
}
