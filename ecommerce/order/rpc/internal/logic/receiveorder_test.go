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

func TestReceiveOrderLogic_ReceiveOrder(t *testing.T) {
	configFile := flag.String("f", "../../etc/order.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test orders
	shippedOrder := &model.Orders{
		OrderNo:     "TEST_ORDER_001",
		UserId:      1,
		TotalAmount: 100,
		PayAmount:   100,
		Status:      3, // Shipped
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	unshippedOrder := &model.Orders{
		OrderNo:     "TEST_ORDER_002",
		UserId:      1,
		TotalAmount: 100,
		PayAmount:   100,
		Status:      2, // Unshipped
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	result1, _ := ctx.OrdersModel.Insert(context.Background(), shippedOrder)
	result2, _ := ctx.OrdersModel.Insert(context.Background(), unshippedOrder)
	shippedOrderId, _ := result1.LastInsertId()
	unshippedOrderId, _ := result2.LastInsertId()

	// Create shipping record
	shipping := &model.OrderShipping{
		OrderId: uint64(shippedOrderId),
		ShippingNo: sql.NullString{
			String: "SF123456",
			Valid:  true,
		},
		Company: sql.NullString{
			String: "SF Express",
			Valid:  true,
		},
		Status: 1, // Shipped
		ShipTime: sql.NullTime{
			Time:  time.Now(),
			Valid: true,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
	ctx.OrderShippingModel.Insert(context.Background(), shipping)

	tests := []struct {
		name    string
		req     *order.ReceiveOrderRequest
		wantErr error
	}{
		{
			name: "receive order successfully",
			req: &order.ReceiveOrderRequest{
				OrderNo: "TEST_ORDER_001",
			},
			wantErr: nil,
		},
		{
			name: "empty order number",
			req: &order.ReceiveOrderRequest{
				OrderNo: "",
			},
			wantErr: zeroerr.ErrOrderNoEmpty,
		},
		{
			name: "order not found",
			req: &order.ReceiveOrderRequest{
				OrderNo: "NOT_EXIST_ORDER",
			},
			wantErr: model.ErrNotFound,
		},
		{
			name: "invalid order status",
			req: &order.ReceiveOrderRequest{
				OrderNo: "TEST_ORDER_002",
			},
			wantErr: zeroerr.ErrOrderStatusNotAllowed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewReceiveOrderLogic(context.Background(), ctx)
			resp, err := l.ReceiveOrder(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)
			}
		})
	}

	// Cleanup
	_ = ctx.OrdersModel.Delete(context.Background(), uint64(shippedOrderId))
	_ = ctx.OrdersModel.Delete(context.Background(), uint64(unshippedOrderId))
}
