package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/producer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/google/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type ShipOrderLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewShipOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ShipOrderLogic {
	return &ShipOrderLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ShipOrderLogic) ShipOrder(in *order.ShipOrderRequest) (*order.ShipOrderResponse, error) {
	if len(in.OrderNo) == 0 || len(in.ShippingNo) == 0 || len(in.Company) == 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	orderInfo, err := l.svcCtx.OrdersModel.FindOneByOrderNo(l.ctx, in.OrderNo)
	if err != nil {
		return nil, err
	}

	if orderInfo.Status != 2 { 
		return nil, zeroerr.ErrOrderStatusNotAllowed
	}

	// Update order status to shipped
	err = l.svcCtx.OrdersModel.UpdateStatus(l.ctx, orderInfo.Id, 3)
	if err != nil {
		return nil, err
	}

	shipping := &model.OrderShipping{
		OrderId: orderInfo.Id,
		ShippingNo: sql.NullString{
			String: in.ShippingNo,
			Valid:  in.ShippingNo != "",
		},
		Company: sql.NullString{
			String: in.Company,
			Valid:  in.Company != "",
		},
		Status: 1, 
		ShipTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = l.svcCtx.OrderShippingModel.Insert(l.ctx, shipping)
	if err != nil {
		return nil, err
	}

	// After shipping record creation, publish event
    event := producer.CreateOrderEvent(
        uuid.New().String(),
        types.EventTypeOrderShipped,
        &types.OrderShippedData{
            OrderNo:    orderInfo.OrderNo,
            OrderId:    int64(orderInfo.Id),
            ShippingNo: in.ShippingNo,
            Company:    in.Company,
        },
        types.Metadata{
            Source:  "order.service",
            UserID:  int64(orderInfo.UserId),
            TraceID: l.ctx.Value("trace_id").(string),
        },
    )

    if err := l.svcCtx.Producer.PublishEventSync(event); err != nil {
        return nil, fmt.Errorf("failed to publish shipping event: %w", err)
    }

	return &order.ShipOrderResponse{
		Success: true,
	}, nil
}
