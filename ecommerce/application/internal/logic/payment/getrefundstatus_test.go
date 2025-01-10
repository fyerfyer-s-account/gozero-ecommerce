package payment

import (
    "context"
    "errors"
    "testing"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rpc/payment"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
)

func TestGetRefundStatusLogic_GetRefundStatus(t *testing.T) {
    tests := []struct {
        name    string
        req     *types.RefundStatusReq
        mock    func(mockPayment *Payment)
        want    *types.RefundStatusResp
        wantErr error
    }{
        {
            name: "successful status retrieval",
            req: &types.RefundStatusReq{
                RefundNo: "TEST_REF_001",
            },
            mock: func(mockPayment *Payment) {
                refundTime := time.Now().Unix()
                mockPayment.EXPECT().GetRefund(
                    mock.Anything,
                    &payment.GetRefundRequest{
                        RefundNo: "TEST_REF_001",
                    },
                ).Return(&payment.GetRefundResponse{
                    Refund: &payment.RefundOrder{
                        RefundNo:   "TEST_REF_001",
                        Status:     3,
                        Amount:     50.00,
                        Reason:     "Product quality issue",
                        RefundTime: refundTime,
                    },
                }, nil)
            },
            want: &types.RefundStatusResp{
                Status:     3,
                Amount:     50.00,
                Reason:     "Product quality issue",
                RefundTime: time.Now().Unix(),
            },
            wantErr: nil,
        },
        {
            name: "empty refund number",
            req: &types.RefundStatusReq{
                RefundNo: "",
            },
            mock:    func(mockPayment *Payment) {},
            want:    nil,
            wantErr: zeroerr.ErrInvalidParameter,
        },
        {
            name: "refund not found",
            req: &types.RefundStatusReq{
                RefundNo: "NON_EXISTENT",
            },
            mock: func(mockPayment *Payment) {
                mockPayment.EXPECT().GetRefund(
                    mock.Anything,
                    &payment.GetRefundRequest{
                        RefundNo: "NON_EXISTENT",
                    },
                ).Return(nil, errors.New("refund not found"))
            },
            want:    nil,
            wantErr: errors.New("refund not found"),
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockPayment := NewPayment(t)
            tt.mock(mockPayment)

            svcCtx := &svc.ServiceContext{
                PaymentRpc: mockPayment,
            }

            logic := NewGetRefundStatusLogic(context.Background(), svcCtx)
            got, err := logic.GetRefundStatus(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr.Error(), err.Error())
            } else {
                assert.NoError(t, err)
                assert.Equal(t, tt.want.Status, got.Status)
                assert.Equal(t, tt.want.Amount, got.Amount)
                assert.Equal(t, tt.want.Reason, got.Reason)
                assert.NotZero(t, got.RefundTime)
            }
        })
    }
}