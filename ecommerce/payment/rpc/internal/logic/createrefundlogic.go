package logic

import (
	"context"
    "database/sql"
    "encoding/json"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/util"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
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

// 退款相关
func (l *CreateRefundLogic) CreateRefund(in *payment.CreateRefundRequest) (*payment.CreateRefundResponse, error) {
    // Input validation
    if in.PaymentNo == "" || in.OrderNo == "" || in.Amount <= 0 || in.Reason == "" {
        return nil, zeroerr.ErrInvalidParam
    }

    // Get payment order
    paymentOrder, err := l.svcCtx.PaymentOrdersModel.FindOneByPaymentNo(l.ctx, in.PaymentNo)
    if err != nil {
        return nil, zeroerr.ErrPaymentNotFound
    }

    // Validate payment status
    if paymentOrder.Status != 3 { // 3: Paid
        return nil, zeroerr.ErrInvalidPaymentStatus
    }

    // Verify refund amount
    if in.Amount > paymentOrder.Amount {
        return nil, zeroerr.ErrRefundAmountInvalid
    }

    // Check existing refunds
    existingRefunds, err := l.svcCtx.RefundOrdersModel.FindByPaymentNo(l.ctx, in.PaymentNo)
    if err != nil {
        return nil, err
    }

    var totalRefunded float64
    for _, refund := range existingRefunds {
        if refund.Status == 3 { // 3: Refunded
            totalRefunded += refund.Amount
        }
    }

    if totalRefunded+in.Amount > paymentOrder.Amount {
        return nil, zeroerr.ErrRefundExceedAmount
    }

    // Generate refund number
    refundNo := util.GenerateNo("REF")

    // Create refund order
    refund := &model.RefundOrders{
        RefundNo:  refundNo,
        PaymentNo: in.PaymentNo,
        OrderNo:   in.OrderNo,
        UserId:    paymentOrder.UserId,
        Amount:    in.Amount,
        Reason:    in.Reason,
        Status:    1, // Pending
        NotifyUrl: sql.NullString{String: in.NotifyUrl, Valid: in.NotifyUrl != ""},
    }

    _, err = l.svcCtx.RefundOrdersModel.Insert(l.ctx, refund)
    if err != nil {
        return nil, err
    }

    // Log refund request
    logData := map[string]interface{}{
        "refund_no":  refundNo,
        "amount":     in.Amount,
        "reason":     in.Reason,
        "notify_url": in.NotifyUrl,
    }
    requestData, _ := json.Marshal(logData)

    log := &model.PaymentLogs{
        PaymentNo: in.PaymentNo,
        Type:      2, // Refund
        Channel:   paymentOrder.Channel,
        RequestData: sql.NullString{
            String: string(requestData),
            Valid:  true,
        },
    }

    _, err = l.svcCtx.PaymentLogsModel.Insert(l.ctx, log)
    if err != nil {
        logx.Errorf("Failed to create refund log: %v", err)
    }

    return &payment.CreateRefundResponse{
        RefundNo: refundNo,
    }, nil
}
