package payment

import (
    "context"
    "errors"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestCreatePaymentLogic_CreatePayment(t *testing.T) {
    tests := []struct {
        name    string
        userId  int64
        req     *types.CreatePaymentReq
        mock    func(mockPayment *Payment)
        want    *types.CreatePaymentResp
        wantErr error
    }{
        {
            name:   "successful payment creation",
            userId: 12345,
            req: &types.CreatePaymentReq{
                OrderNo:     "TEST_ORDER_001",
                PaymentType: 1, // WeChat
                Amount:      100.00,
                NotifyUrl:   "http://example.com/notify",
                ReturnUrl:   "http://example.com/return",
            },
            mock: func(mockPayment *Payment) {
                mockPayment.EXPECT().CreatePayment(
                    mock.Anything,
                    &payment.CreatePaymentRequest{
                        OrderNo:   "TEST_ORDER_001",
                        UserId:    12345,
                        Amount:    100.00,
                        Channel:   1,
                        NotifyUrl: "http://example.com/notify",
                        ReturnUrl: "http://example.com/return",
                    },
                ).Return(&payment.CreatePaymentResponse{
                    PaymentNo: "PAY_TEST_001",
                    PayUrl:    "weixin://pay?orderNo=PAY_TEST_001",
                }, nil)
            },
            want: &types.CreatePaymentResp{
                PaymentNo: "PAY_TEST_001",
                PayUrl:    "weixin://pay?orderNo=PAY_TEST_001",
            },
            wantErr: nil,
        },
        {
            name:   "rpc error",
            userId: 12345,
            req: &types.CreatePaymentReq{
                OrderNo:     "TEST_ORDER_002",
                PaymentType: 2,
                Amount:      200.00,
            },
            mock: func(mockPayment *Payment) {
                mockPayment.EXPECT().CreatePayment(
                    mock.Anything,
                    &payment.CreatePaymentRequest{
                        OrderNo: "TEST_ORDER_002",
                        UserId:  12345,
                        Amount:  200.00,
                        Channel: 2,
                    },
                ).Return(nil, errors.New("rpc error"))
            },
            want:    nil,
            wantErr: errors.New("rpc error"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Create mock payment client
            mockPayment := NewPayment(t)

            // Setup mock expectations
            tt.mock(mockPayment)

            // Create service context with mock
            svcCtx := &svc.ServiceContext{
                PaymentRpc: mockPayment,
            }

            // Create context with userId
            ctx := context.WithValue(context.Background(), "userId", tt.userId)

            // Create logic instance
            logic := NewCreatePaymentLogic(ctx, svcCtx)

            // Execute CreatePayment
            got, err := logic.CreatePayment(tt.req)

            // Assert results
            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr.Error(), err.Error())
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want, got)
            }
        })
    }
}