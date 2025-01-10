package logic

import (
    "context"
    "fmt"
    "time"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/util"
    "github.com/zeromicro/go-zero/core/logx"
	"database/sql"
)

type CreatePaymentLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewCreatePaymentLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreatePaymentLogic {
    return &CreatePaymentLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *CreatePaymentLogic) CreatePayment(in *payment.CreatePaymentRequest) (*payment.CreatePaymentResponse, error) {
    // Input validation
    if in.OrderNo == "" || in.UserId == 0 || in.Channel == 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    if in.Amount < l.svcCtx.Config.MinAmount || in.Amount > l.svcCtx.Config.MaxAmount {
        return nil, zeroerr.ErrInvalidAmount
    }

    // Check payment channel
    _, err := l.svcCtx.PaymentChannelsModel.FindOneByChannelAndStatus(l.ctx, in.Channel, 1)
    if err != nil {
        return nil, zeroerr.ErrChannelNotFound
    }

    // Generate payment number
    paymentNo := util.GenerateNo("PAY")

    // Create payment order
    expireTime := time.Now().Add(time.Duration(l.svcCtx.Config.PaymentTimeout) * time.Second)
    p := &model.PaymentOrders{
        PaymentNo: paymentNo,
        OrderNo:   in.OrderNo,
        UserId:    uint64(in.UserId),
        Amount:    in.Amount,
        Channel:   in.Channel,
        Status:    1, // Pending payment
        NotifyUrl: sql.NullString{String: in.NotifyUrl, Valid: in.NotifyUrl != ""},
        ReturnUrl: sql.NullString{String: in.ReturnUrl, Valid: in.ReturnUrl != ""},
        ExpireTime: sql.NullTime{
            Time:  expireTime,
            Valid: true,
        },
    }

    _, err = l.svcCtx.PaymentOrdersModel.Insert(l.ctx, p)
    if err != nil {
        return nil, err
    }

    // Generate payment URL/params based on channel
    var payUrl string
    switch in.Channel {
    case 1: // WeChat
        payUrl = fmt.Sprintf("wechat://pay?orderNo=%s&amount=%.2f", paymentNo, in.Amount)
    case 2: // Alipay
        payUrl = fmt.Sprintf("alipay://pay?orderNo=%s&amount=%.2f", paymentNo, in.Amount)
    case 3: // Balance
        payUrl = fmt.Sprintf("balance://pay?orderNo=%s&amount=%.2f", paymentNo, in.Amount)
    }

    return &payment.CreatePaymentResponse{
        PaymentNo: paymentNo,
        PayUrl:    payUrl,
    }, nil
}