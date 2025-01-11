package logic

import (
	"context"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetRefundLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetRefundLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetRefundLogic {
	return &GetRefundLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetRefundLogic) GetRefund(in *order.GetRefundRequest) (*order.GetRefundResponse, error) {
    if len(in.RefundNo) == 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    refund, err := l.svcCtx.OrderRefundsModel.FindOneByRefundNo(l.ctx, in.RefundNo)
    if err != nil {
        return nil, err
    }

    // Parse images JSON string to array
    var images []string
    if refund.Images.Valid && len(refund.Images.String) > 0 {
        if err := json.Unmarshal([]byte(refund.Images.String), &images); err != nil {
            return nil, err
        }
    }

    return &order.GetRefundResponse{
        Refund: &order.RefundInfo{
            Id:          int64(refund.Id),
            OrderId:     int64(refund.OrderId),
            RefundNo:    refund.RefundNo,
            Amount:      refund.Amount,
            Reason:      refund.Reason,
            Status:      int32(refund.Status),
            Description: refund.Description.String,
            Images:      images,
            CreatedAt:   refund.CreatedAt.Unix(),
            UpdatedAt:   refund.UpdatedAt.Unix(),
        },
    }, nil
}
