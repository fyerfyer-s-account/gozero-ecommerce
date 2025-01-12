package order

import (
	"context"
	"math"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/orderservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListOrdersLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListOrdersLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListOrdersLogic {
	return &ListOrdersLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListOrdersLogic) ListOrders(req *types.OrderListReq) (*types.OrderListResp, error) {
	resp, err := l.svcCtx.OrderRpc.ListOrders(l.ctx, &orderservice.ListOrdersRequest{
		UserId:   l.ctx.Value("userId").(int64),
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
	})
	if err != nil {
		return nil, err
	}

	orders := make([]types.Order, 0)
	for _, o := range resp.Orders {
		orders = append(orders, *convertTypeOrder(o))
	}

	var totalPages int32
	if resp.Total > 0 {
		totalPages = int32(math.Ceil(float64(resp.Total) / float64(req.PageSize)))
	}

	return &types.OrderListResp{
		List:       orders,
		Total:      resp.Total,
		Page:       req.Page,
		TotalPages: totalPages,
	}, nil
}
