package logic

import (
    "context"
    "database/sql"
    "encoding/json"
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

func TestGetRefundLogic_GetRefund(t *testing.T) {
    configFile := flag.String("f", "../../etc/order.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    images := []string{"image1.jpg", "image2.jpg"}
    imagesJSON, _ := json.Marshal(images)

    // Create test refund
    testRefund := &model.OrderRefunds{
        OrderId:     1,
        RefundNo:    "RF_TEST_001",
        Amount:      50.0,
        Reason:      "test refund",
        Status:      0,
        Description: sql.NullString{String: "test description", Valid: true},
        Images:      sql.NullString{String: string(imagesJSON), Valid: true},
        CreatedAt:   time.Now(),
        UpdatedAt:   time.Now(),
    }
    
    result, _ := ctx.OrderRefundsModel.Insert(context.Background(), testRefund)
    refundId, _ := result.LastInsertId()

    tests := []struct {
        name    string
        req     *order.GetRefundRequest
        wantErr error
    }{
        {
            name: "get refund successfully",
            req: &order.GetRefundRequest{
                RefundNo: "RF_TEST_001",
            },
            wantErr: nil,
        },
        {
            name: "empty refund number",
            req: &order.GetRefundRequest{
                RefundNo: "",
            },
            wantErr: zeroerr.ErrInvalidParam,
        },
        {
            name: "refund not found",
            req: &order.GetRefundRequest{
                RefundNo: "NOT_EXIST_REFUND",
            },
            wantErr: model.ErrNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            l := NewGetRefundLogic(context.Background(), ctx)
            resp, err := l.GetRefund(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.Equal(t, "RF_TEST_001", resp.Refund.RefundNo)
                assert.Equal(t, float64(50.0), resp.Refund.Amount)
                assert.Equal(t, "test description", resp.Refund.Description)
                assert.Equal(t, images, resp.Refund.Images)
            }
        })
    }

    // Cleanup
    _ = ctx.OrderRefundsModel.Delete(context.Background(), uint64(refundId))
}