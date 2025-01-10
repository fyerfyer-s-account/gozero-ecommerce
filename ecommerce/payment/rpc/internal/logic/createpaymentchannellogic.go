package logic

import (
	"context"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/zeromicro/go-zero/core/logx"
)

type CreatePaymentChannelLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewCreatePaymentChannelLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentChannelLogic {
    return &CreatePaymentChannelLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *CreatePaymentChannelLogic) CreatePaymentChannel(req *payment.CreatePaymentChannelRequest) (*payment.CreatePaymentChannelResponse, error) {
    // Input validation
    if req.Name == "" || req.Channel == 0 || req.Config == "" {
        return nil, zeroerr.ErrInvalidParam
    }

    // Check if channel already exists
    _, err := l.svcCtx.PaymentChannelsModel.FindOneByChannel(l.ctx, int64(req.Channel))
    if err == nil {
        return nil, zeroerr.ErrPaymentChannelExists
    }

    // Validate config JSON
    var configMap map[string]interface{}
    if err := json.Unmarshal([]byte(req.Config), &configMap); err != nil {
        return nil, zeroerr.ErrInvalidParam
    }

    // Create payment channel
    channel := &model.PaymentChannels{
        Name:    req.Name,
        Channel: int64(req.Channel),
        Config:  req.Config,
        Status:  1, 
    }

    result, err := l.svcCtx.PaymentChannelsModel.Insert(l.ctx, channel)
    if err != nil {
        return nil, err
    }

    channelId, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }

    return &payment.CreatePaymentChannelResponse{
        Id: int64(channelId),
    }, nil
}