package logic

import (
	"context"
	"database/sql"
	"flag"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/order"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestGetPaymentStatusLogic_GetPaymentStatus(t *testing.T) {
	configFile := flag.String("f", "../../etc/order.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test payment
	testPayment := &model.OrderPayments{
		OrderId:       1,
		PaymentNo:     "PAY_TEST_001",
		PaymentMethod: 1,
		Amount:        100,
		Status:        1,
		PayTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	result, _ := ctx.OrderPaymentsModel.Insert(context.Background(), testPayment)
	paymentId, _ := result.LastInsertId()

	tests := []struct {
		name    string
		req     *order.GetPaymentStatusRequest
		want    int32
		wantErr error
	}{
		{
			name: "get payment status successfully",
			req: &order.GetPaymentStatusRequest{
				PaymentNo: "PAY_TEST_001",
			},
			want:    1,
			wantErr: nil,
		},
		{
			name: "empty payment number",
			req: &order.GetPaymentStatusRequest{
				PaymentNo: "",
			},
			want:    0,
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "payment not found",
			req: &order.GetPaymentStatusRequest{
				PaymentNo: "NOT_EXIST_PAYMENT",
			},
			want:    0,
			wantErr: model.ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewGetPaymentStatusLogic(context.Background(), ctx)
			resp, err := l.GetPaymentStatus(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.Equal(t, tt.want, resp.Status)
			}
		})
	}

	// Cleanup
	_ = ctx.OrderPaymentsModel.Delete(context.Background(), uint64(paymentId))
}
