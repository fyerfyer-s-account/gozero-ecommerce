package logic

import (
    "context"
    "flag"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
    "github.com/stretchr/testify/assert"
    "github.com/zeromicro/go-zero/core/conf"
)

func TestCreateStockInLogic(t *testing.T) {
    configFile := flag.String("f", "../../etc/inventory.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    logic := NewCreateStockInLogic(context.Background(), ctx)

    tests := []struct {
        name    string
        req     *inventory.CreateStockInRequest
        wantErr bool
    }{
        {
            name: "normal case",
            req: &inventory.CreateStockInRequest{
                WarehouseId: 1,
                Items: []*inventory.StockInItem{
                    {
                        SkuId:    1001,
                        Quantity: 100,
                    },
                    {
                        SkuId:    1002,
                        Quantity: 200,
                    },
                },
                Remark: "Test stock in",
            },
            wantErr: false,
        },
        {
            name: "invalid warehouse id",
            req: &inventory.CreateStockInRequest{
                WarehouseId: 0,
                Items: []*inventory.StockInItem{
                    {
                        SkuId:    1001,
                        Quantity: 100,
                    },
                },
            },
            wantErr: true,
        },
        {
            name: "empty items",
            req: &inventory.CreateStockInRequest{
                WarehouseId: 1,
                Items:      []*inventory.StockInItem{},
            },
            wantErr: true,
        },
        {
            name: "invalid quantity",
            req: &inventory.CreateStockInRequest{
                WarehouseId: 1,
                Items: []*inventory.StockInItem{
                    {
                        SkuId:    1001,
                        Quantity: -1,
                    },
                },
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            resp, err := logic.CreateStockIn(tt.req)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.True(t, resp.Success)

                // Verify stock quantities
                for _, item := range tt.req.Items {
                    stock, err := ctx.StocksModel.FindOneBySkuIdWarehouseId(context.Background(), 
                        uint64(item.SkuId), uint64(tt.req.WarehouseId))
                    assert.NoError(t, err)
                    assert.GreaterOrEqual(t, stock.Available, int64(item.Quantity))
                }
            }
        })
    }
}