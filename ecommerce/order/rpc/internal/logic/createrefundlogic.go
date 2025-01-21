package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
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

type CreateRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateRefundLogic {
	return &CreateRefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *CreateRefundLogic) CreateRefund(in *order.CreateRefundRequest) (*order.CreateRefundResponse, error) {
    if len(in.OrderNo) == 0 || in.Amount <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    orderInfo, err := l.svcCtx.OrdersModel.FindOneByOrderNo(l.ctx, in.OrderNo)
    if err != nil {
        return nil, err
    }

    if orderInfo.Status != 2 && orderInfo.Status != 3 && orderInfo.Status != 4 {
        return nil, zeroerr.ErrRefundNotAllowed
    }

    if in.Amount > orderInfo.PayAmount {
        return nil, zeroerr.ErrRefundExceedAmount
    }

    refundNo := fmt.Sprintf("RF%d%d", time.Now().UnixNano(), orderInfo.UserId)

    // Convert images array to JSON string
    var imagesJSON string
    if len(in.Images) > 0 {
        imagesBytes, err := json.Marshal(in.Images)
        if err != nil {
            return nil, err
        }
        imagesJSON = string(imagesBytes)
    }

    refund := &model.OrderRefunds{
        OrderId:     orderInfo.Id,
        RefundNo:    refundNo,
        Amount:      in.Amount,
        Reason:      in.Reason,
        Status:      0,
        Description: sql.NullString{String: in.Description, Valid: len(in.Description) > 0},
        Images:      sql.NullString{String: imagesJSON, Valid: len(imagesJSON) > 0},
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }

    _, err = l.svcCtx.OrderRefundsModel.Insert(l.ctx, refund)
    if err != nil {
        return nil, err
    }

    err = l.svcCtx.OrdersModel.UpdateStatus(l.ctx, orderInfo.Id, 6)
    if err != nil {
        return nil, err
    }

    event := producer.CreateOrderEvent(
        uuid.New().String(),
        types.EventTypeRefundCreated,
        &types.RefundCreatedData{
            OrderNo:  orderInfo.OrderNo,
            OrderId:  int64(orderInfo.Id),
            RefundNo: refundNo,
            Amount:   in.Amount,
            Reason:   in.Reason,
        },
        types.Metadata{
            Source:  "order.service",
            UserID:  int64(orderInfo.UserId),
            TraceID: l.ctx.Value("trace_id").(string),
        },
    )

    if err := l.svcCtx.Producer.PublishEventSync(event); err != nil {
        return nil, fmt.Errorf("failed to publish refund created event: %w", err)
    }

    return &order.CreateRefundResponse{
        RefundNo: refundNo,
    }, nil
}

func joinImages(images []string) string {
    if len(images) == 0 {
        return ""
    }
    return strings.Join(images, ",")
}
