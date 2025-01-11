package logic

import (
    "context"
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

func TestListOrdersLogic_ListOrders(t *testing.T) {
    configFile := flag.String("f", "../../etc/order.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    // Create test orders
    orders := []*model.Orders{
        {
            OrderNo:      "TEST_ORDER_001",
            UserId:      1,
            TotalAmount: 100,
            PayAmount:   100,
            Status:      1,
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        },
        {
            OrderNo:      "TEST_ORDER_002",
            UserId:      1,
            TotalAmount: 200,
            PayAmount:   200,
            Status:      2,
            CreatedAt:   time.Now(),
            UpdatedAt:   time.Now(),
        },
    }

    var orderIds []uint64
    for _, o := range orders {
        result, _ := ctx.OrdersModel.Insert(context.Background(), o)
        id, _ := result.LastInsertId()
        orderIds = append(orderIds, uint64(id))
    }

    tests := []struct {
        name    string
        req     *order.ListOrdersRequest
        want    int
        wantErr error
    }{
        {
            name: "list orders successfully",
            req: &order.ListOrdersRequest{
                UserId:   1,
                Status:   -1,
                Page:     1,
                PageSize: 10,
            },
            want:    2,
            wantErr: nil,
        },
        {
            name: "filter by status",
            req: &order.ListOrdersRequest{
                UserId:   1,
                Status:   1,
                Page:     1,
                PageSize: 10,
            },
            want:    1,
            wantErr: nil,
        },
        {
            name: "invalid user id",
            req: &order.ListOrdersRequest{
                UserId:   0,
                Page:     1,
                PageSize: 10,
            },
            want:    0,
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "invalid page",
            req: &order.ListOrdersRequest{
                UserId:   1,
                Page:     0,
                PageSize: 10,
            },
            want:    0,
            wantErr: zeroerr.ErrInvalidParam,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewListOrdersLogic(context.Background(), ctx)
            resp, err := l.ListOrders(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.Equal(t, tt.want, len(resp.Orders))
                if tt.want > 0 {
                    assert.Equal(t, int64(1), resp.Orders[0].UserId)
                }
            }
        })
    }

    // Cleanup
    for _, id := range orderIds {
        _ = ctx.OrdersModel.Delete(context.Background(), id)
    }
}