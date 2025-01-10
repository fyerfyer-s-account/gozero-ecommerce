package payment

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestGetPaymentStatusLogic_GetPaymentStatus(t *testing.T) {
    tests := []struct {
        name    string
        req     *types.PaymentStatusReq
        mock    func(mockPayment *Payment)
        want    *types.PaymentStatusResp
        wantErr error
    }{
        {
            name: "successful status retrieval",
            req: &types.PaymentStatusReq{
                PaymentNo: "TEST_PAY_001",
            },
            mock: func(mockPayment *Payment) {
                payTime := time.Now().Unix()
                mockPayment.EXPECT().GetPayment(
                    mock.Anything,
                    &payment.GetPaymentRequest{
                        PaymentNo: "TEST_PAY_001",
                    },
                ).Return(&payment.GetPaymentResponse{
                    Payment: &payment.PaymentOrder{
                        PaymentNo: "TEST_PAY_001",
                        Status:    3, // Paid
                        Amount:    100.00,
                        PayTime:   payTime,
                    },
                }, nil)
            },
            want: &types.PaymentStatusResp{
                Status:   3,
                Amount:   100.00,
                PayTime:  time.Now().Unix(),
                ErrorMsg: "",
            },
            wantErr: nil,
        },
        {
            name: "empty payment number",
            req: &types.PaymentStatusReq{
                PaymentNo: "",
            },
            mock:    func(mockPayment *Payment) {},
            want:    nil,
            wantErr: errors.New("payment number cannot be empty"),
        },
        {
            name: "payment not found",
            req: &types.PaymentStatusReq{
                PaymentNo: "NON_EXISTENT",
            },
            mock: func(mockPayment *Payment) {
                mockPayment.EXPECT().GetPayment(
                    mock.Anything,
                    &payment.GetPaymentRequest{
                        PaymentNo: "NON_EXISTENT",
                    },
                ).Return(nil, errors.New("payment not found"))
            },
            want:    nil,
            wantErr: errors.New("payment not found"),
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

            // Create logic instance
            logic := NewGetPaymentStatusLogic(context.Background(), svcCtx)

            // Execute GetPaymentStatus
            got, err := logic.GetPaymentStatus(tt.req)

            // Assert results
            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr.Error(), err.Error())
            } else {
                assert.NoError(t, err)
                // Compare fields individually since PayTime might be different
                assert.Equal(t, tt.want.Status, got.Status)
                assert.Equal(t, tt.want.Amount, got.Amount)
                assert.NotZero(t, got.PayTime)
                assert.Equal(t, tt.want.ErrorMsg, got.ErrorMsg)
            }
        })
    }
}