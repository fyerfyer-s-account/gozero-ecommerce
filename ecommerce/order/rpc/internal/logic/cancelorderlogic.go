package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/producer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/google/uuid"

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

    event := producer.CreateOrderEvent(
        uuid.New().String(),
        types.EventTypeOrderCancelled,
        &types.OrderCancelledData{
            OrderNo: orderInfo.OrderNo,
            OrderId: int64(orderInfo.Id),
            Amount:  orderInfo.PayAmount,
            Reason:  in.Reason,
        },
        types.Metadata{
            Source:  "order.service",
            UserID:  int64(orderInfo.UserId),
            TraceID: l.ctx.Value("trace_id").(string),
        },
    )

    if err := l.svcCtx.Producer.PublishEventSync(event); err != nil {
        return nil, err
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