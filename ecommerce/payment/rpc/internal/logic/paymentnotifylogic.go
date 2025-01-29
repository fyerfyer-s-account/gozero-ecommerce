package logic

import (
	"context"
	"database/sql"
	"encoding/json"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type PaymentNotifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewPaymentNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *PaymentNotifyLogic {
	return &PaymentNotifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *PaymentNotifyLogic) PaymentNotify(in *payment.PaymentNotifyRequest) (*payment.PaymentNotifyResponse, error) {
    if in.Channel == 0 || in.NotifyData == "" {
        return nil, zeroerr.ErrInvalidParam
    }

    // Parse notify data based on channel
    var paymentNo string
    var success bool
    var err error

    switch in.Channel {
    case 1: // WeChat
        paymentNo, success, err = parseWechatNotify(in.NotifyData)
    case 2: // Alipay
        paymentNo, success, err = parseAlipayNotify(in.NotifyData)
    case 3: // Balance
        paymentNo, success, err = parseBalanceNotify(in.NotifyData)
    default:
        return nil, zeroerr.ErrUnsupportedChannel
    }

    if err != nil {
        return nil, zeroerr.ErrInvalidNotifyData
    }

    // Log notification
    log := &model.PaymentLogs{
        PaymentNo: paymentNo,
        Type:      1, // Payment
        Channel:   int64(in.Channel),
        RequestData: sql.NullString{
            String: in.NotifyData,
            Valid:  true,
        },
    }

    _, err = l.svcCtx.PaymentLogsModel.Insert(l.ctx, log)
    if err != nil {
        logx.Errorf("Failed to log payment notification: %v", err)
    }

    // Update payment status
    p, err := l.svcCtx.PaymentOrdersModel.FindOneByPaymentNo(l.ctx, paymentNo)
    if err != nil {
        return nil, zeroerr.ErrPaymentNotFound
    }

    if success {
        err = l.svcCtx.PaymentOrdersModel.UpdatePartial(l.ctx, p.Id, map[string]interface{}{
            "status": 3, // Paid
            "pay_time": time.Now(),
            "channel_data": in.NotifyData,
        })

        if err != nil {
            return nil, err
        }

        // Publish payment success event
        successEvent := &types.PaymentSuccessEvent {
            Amount: p.Amount,
            PaymentMethod: in.Channel,
            PaidTime: time.Now(),
        }

        if err := l.svcCtx.Producer.PublishPaymentSuccess(l.ctx, successEvent); err != nil {
            logx.Errorf("Failed to publish payment success event: %v", err)
        }

    } else {
        err = l.svcCtx.PaymentOrdersModel.UpdateStatus(l.ctx, p.Id, 5) // Closed

        if err != nil {
            return nil, err 
        }

        // Publish payment failed event
        failedEvent := &types.PaymentFailedEvent {
            Amount: p.Amount,
            Reason: "Payment notification failed",
            ErrorCode: "PAYMENT_NOTIFY_FAILED",
        }

        if err := l.svcCtx.Producer.PublishPaymentFailed(l.ctx, failedEvent); err != nil {
            logx.Errorf("Failed to publish payment failed event: %v", err)
        }
    }

    if err != nil {
        return nil, err
    }

    return &payment.PaymentNotifyResponse{
        ReturnMsg: "success",
    }, nil
}

func parseWechatNotify(data string) (string, bool, error) {
    var notify struct {
        OutTradeNo string `json:"out_trade_no"`
        ResultCode string `json:"result_code"`
    }
    if err := json.Unmarshal([]byte(data), &notify); err != nil {
        return "", false, err
    }
    return notify.OutTradeNo, notify.ResultCode == "SUCCESS", nil
}

func parseAlipayNotify(data string) (string, bool, error) {
    var notify struct {
        OutTradeNo string `json:"out_trade_no"`
        TradeStatus string `json:"trade_status"`
    }
    if err := json.Unmarshal([]byte(data), &notify); err != nil {
        return "", false, err
    }
    return notify.OutTradeNo, notify.TradeStatus == "TRADE_SUCCESS", nil
}

func parseBalanceNotify(data string) (string, bool, error) {
    var notify struct {
        PaymentNo string `json:"payment_no"`
        Success   bool   `json:"success"`
    }
    if err := json.Unmarshal([]byte(data), &notify); err != nil {
        return "", false, err
    }
    return notify.PaymentNo, notify.Success, nil
}
