package logic

import (
	"context"
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

    event := producer.CreateOrderEvent(
        uuid.New().String(),
        types.EventTypeOrderPaid,
        &types.OrderPaidData{
            OrderNo:       orderInfo.OrderNo,
            PaymentNo:     paymentNo,
            PayAmount:     orderInfo.PayAmount,
            PaymentMethod: int32(in.PaymentMethod),
            PayTime:      time.Now(),
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

    return &order.PayOrderResponse{
        PaymentNo: paymentNo,
        PayUrl:    payURL,
    }, nil
}
