package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/producer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/google/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type ProcessRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewProcessRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ProcessRefundLogic {
	return &ProcessRefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ProcessRefundLogic) ProcessRefund(in *order.ProcessRefundRequest) (*order.ProcessRefundResponse, error) {
	if len(in.RefundNo) == 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	refund, err := l.svcCtx.OrderRefundsModel.FindOneByRefundNo(l.ctx, in.RefundNo)
	if err != nil {
		return nil, err
	}

    o, err := l.svcCtx.OrdersModel.FindOne(l.ctx, refund.OrderId)
    if err != nil {
        return nil, err 
    }

	if refund.Status != 0 {
		return nil, zeroerr.ErrRefundStatusInvalid
	}

	newStatus := int64(2) // Rejected
	if in.Agree {
		newStatus = 1 // Approved
	}

	err = l.svcCtx.OrderRefundsModel.UpdateStatus(l.ctx, in.RefundNo, newStatus, in.Reply)
	if err != nil {
		return nil, err
	}

	if in.Agree {
		// Update order status to refunded
		err = l.svcCtx.OrdersModel.UpdateStatus(l.ctx, refund.OrderId, 7) // 7: Refunded
		if err != nil {
			return nil, err
		}

		// Update payment status
		payment, err := l.svcCtx.OrderPaymentsModel.FindByOrderId(l.ctx, refund.OrderId)
		if err != nil {
			return nil, err
		}

		err = l.svcCtx.OrderPaymentsModel.UpdateStatus(l.ctx, payment.PaymentNo, 2, time.Now()) // 2: Refunded
		if err != nil {
			return nil, err
		}

		// After refund processing, publish event
		event := producer.CreateOrderEvent(
			uuid.New().String(),
			types.EventTypeRefundProcessed,
			&types.RefundProcessedData{
				OrderNo:  o.OrderNo,
				OrderId:  int64(refund.OrderId),
				RefundNo: refund.RefundNo,
				Amount:   refund.Amount,
				Status:   newStatus,
				Reply:    in.Reply,
			},
			types.Metadata{
				Source:  "order.service",
				UserID:  int64(o.UserId),
				TraceID: l.ctx.Value("trace_id").(string),
			},
		)

		if err := l.svcCtx.Producer.PublishEventSync(event); err != nil {
			return nil, fmt.Errorf("failed to publish refund processed event: %w", err)
		}
	}

	return &order.ProcessRefundResponse{
		Success: true,
	}, nil
}
