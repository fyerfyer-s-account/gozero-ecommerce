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

	"github.com/zeromicro/go-zero/core/logx"
)

type RefundNotifyLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRefundNotifyLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RefundNotifyLogic {
	return &RefundNotifyLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RefundNotifyLogic) RefundNotify(in *payment.RefundNotifyRequest) (*payment.RefundNotifyResponse, error) {
    if in.Channel == 0 || in.NotifyData == "" {
        return nil, zeroerr.ErrInvalidParam
    }

    // Parse notify data
    var refundNo string
    var success bool
    var err error

    switch in.Channel {
    case 1: // WeChat
        refundNo, success, err = parseWechatRefundNotify(in.NotifyData)
    case 2: // Alipay
        refundNo, success, err = parseAlipayRefundNotify(in.NotifyData)
    case 3: // Balance
        refundNo, success, err = parseBalanceRefundNotify(in.NotifyData)
    default:
        return nil, zeroerr.ErrUnsupportedChannel
    }

    if err != nil {
        return nil, zeroerr.ErrInvalidNotifyData
    }

    // Get refund order
    refund, err := l.svcCtx.RefundOrdersModel.FindOneByRefundNo(l.ctx, refundNo)
    if err != nil {
        return nil, zeroerr.ErrRefundNotFound
    }

    // Log notification
    log := &model.PaymentLogs{
        PaymentNo: refund.PaymentNo,
        Type:      2, // Refund
        Channel:   int64(in.Channel),
        RequestData: sql.NullString{
            String: in.NotifyData,
            Valid:  true,
        },
    }

    _, err = l.svcCtx.PaymentLogsModel.Insert(l.ctx, log)
    if err != nil {
        logx.Errorf("Failed to log refund notification: %v", err)
    }

    // Update refund status
    updates := map[string]interface{}{
        "channel_data": in.NotifyData,
    }
    if success {
        updates["status"] = 3 // Refunded
        updates["refund_time"] = time.Now()

        // Update payment order status
        payment, err := l.svcCtx.PaymentOrdersModel.FindOneByPaymentNo(l.ctx, refund.PaymentNo)
        if err == nil {
            err = l.svcCtx.PaymentOrdersModel.UpdateStatus(l.ctx, payment.Id, 4) // Refunded
            if err != nil {
                logx.Errorf("Failed to update payment status: %v", err)
            }
        }
    } else {
        updates["status"] = 4 // Refund failed
    }

    err = l.svcCtx.RefundOrdersModel.UpdatePartial(l.ctx, refund.Id, updates)
    if err != nil {
        return nil, err
    }

    return &payment.RefundNotifyResponse{
        ReturnMsg: "success",
    }, nil
}

func parseWechatRefundNotify(data string) (string, bool, error) {
    var notify struct {
        RefundId   string `json:"refund_id"`
        RefundStatus string `json:"refund_status"`
    }
    if err := json.Unmarshal([]byte(data), &notify); err != nil {
        return "", false, err
    }
    return notify.RefundId, notify.RefundStatus == "SUCCESS", nil
}

func parseAlipayRefundNotify(data string) (string, bool, error) {
    var notify struct {
        OutRefundNo string `json:"out_refund_no"`
        RefundStatus string `json:"refund_status"`
    }
    if err := json.Unmarshal([]byte(data), &notify); err != nil {
        return "", false, err
    }
    return notify.OutRefundNo, notify.RefundStatus == "REFUND_SUCCESS", nil
}

func parseBalanceRefundNotify(data string) (string, bool, error) {
    var notify struct {
        RefundNo string `json:"refund_no"`
        Success  bool   `json:"success"`
    }
    if err := json.Unmarshal([]byte(data), &notify); err != nil {
        return "", false, err
    }
    return notify.RefundNo, notify.Success, nil
}
