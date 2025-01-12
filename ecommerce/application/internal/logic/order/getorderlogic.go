package order

import (
	"context"
	"strconv"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetOrderLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetOrderLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetOrderLogic {
	return &GetOrderLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetOrderLogic) GetOrder(req *types.GetOrderReq) (*types.Order, error) {
    orderId := strconv.FormatInt(req.Id, 10)
    resp, err := l.svcCtx.OrderRpc.GetOrder(l.ctx, &order.GetOrderRequest{
        OrderNo: orderId,
    })
    if err != nil {
        return nil, err
    }

    return convertTypeOrder(resp.Order), nil
}

func convertTypeOrder(o *order.Order) *types.Order {
    items := make([]types.OrderItem, 0, len(o.Items))
    for _, item := range o.Items {
        items = append(items, types.OrderItem{
            Id:          item.Id,
            ProductId:   item.ProductId,
            ProductName: item.ProductName,
            SkuId:      item.SkuId,
            SkuName:    item.SkuName,
            Price:      item.Price,
            Quantity:   int32(item.Quantity),
            Amount:     item.TotalAmount,
        })
    }

    var payment types.Payment
    if o.Payment != nil {
        payment = types.Payment{
            PaymentNo:   o.Payment.PaymentNo,
            PaymentType: o.Payment.PaymentMethod,
            Status:      o.Payment.Status,
            Amount:      o.Payment.Amount,
            PayTime:     o.Payment.PayTime,
        }
    }

    var shipping types.Shipping
    if o.Shipping != nil {
        shipping = types.Shipping{
            ShippingNo:  o.Shipping.ShippingNo,
            Company:     o.Shipping.Company,
            Status:      o.Shipping.Status,
            ShipTime:    o.Shipping.ShipTime,
            ReceiveTime: o.Shipping.ReceiveTime,
        }
    }

    return &types.Order{
        Id:            o.Id,
        OrderNo:       o.OrderNo,
        UserId:        o.UserId,
        Status:        o.Status,
        TotalAmount:   o.TotalAmount,
        PayAmount:     o.PayAmount,
        FreightAmount: o.FreightAmount,
        Items:         items,
        Payment:       payment,
        Shipping:      shipping,
        CreatedAt:     o.CreatedAt,
    }
}
