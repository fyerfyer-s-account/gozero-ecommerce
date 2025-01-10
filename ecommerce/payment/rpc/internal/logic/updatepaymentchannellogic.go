package logic

import (
	"context"
	"encoding/json"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdatePaymentChannelLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdatePaymentChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdatePaymentChannelLogic {
	return &UpdatePaymentChannelLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdatePaymentChannelLogic) UpdatePaymentChannel(in *payment.UpdatePaymentChannelRequest) (*payment.UpdatePaymentChannelResponse, error) {
    if in.Id == 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Check if channel exists
    channel, err := l.svcCtx.PaymentChannelsModel.FindOne(l.ctx, uint64(in.Id))
    if err != nil {
        return nil, zeroerr.ErrChannelNotFound
    }

    updates := make(map[string]interface{})

    if in.Name != "" {
        updates["name"] = in.Name
    }

    if in.Config != "" {
        // Validate config JSON
        var configMap map[string]interface{}
        if err := json.Unmarshal([]byte(in.Config), &configMap); err != nil {
            return nil, zeroerr.ErrInvalidChannelConfig
        }
        updates["config"] = in.Config
    }

    if in.Status != 0 {
        if in.Status != 1 && in.Status != 2 {
            return nil, zeroerr.ErrInvalidParam
        }
        updates["status"] = in.Status
    }

    if len(updates) > 0 {
        updates["updated_at"] = time.Now()
        err = l.svcCtx.PaymentChannelsModel.UpdateFields(l.ctx, channel.Id, updates)
        if err != nil {
            return nil, zeroerr.ErrChannelUpdateFailed
        }
    }

    return &payment.UpdatePaymentChannelResponse{
        Success: true,
    }, nil
}
