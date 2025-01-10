package logic

import (
    "context"
    "flag"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
    "github.com/stretchr/testify/assert"
    "github.com/zeromicro/go-zero/core/conf"
)

func TestAddItemLogic_AddItem(t *testing.T) {
    configFile := flag.String("f", "../../etc/cart.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    ctx := svc.NewServiceContext(c)

    mockProduct := NewProductService(t)
    ctx.ProductRpc = mockProduct

    tests := []struct {
        name    string
        req     *cart.AddItemRequest
        mock    func()
        wantErr error
    }{
        {
            name: "Add new item success",
            req: &cart.AddItemRequest{
                UserId:    1,
                ProductId: 1,
                SkuId:    1,
                Quantity: 2,
            },
            mock: func() {
                mockProduct.EXPECT().GetSku(context.Background(), &product.GetSkuRequest{
                    Id: 1,
                }).Return(&product.GetSkuResponse{
                    Sku: &product.Sku{
                        Id:      1,
                        SkuCode: "SKU001",
                        Price:   99.9,
                        Stock:   10,
                    },
                }, nil)

                mockProduct.EXPECT().GetProduct(context.Background(), &product.GetProductRequest{
                    Id: 1,
                }).Return(&product.GetProductResponse{
                    Product: &product.Product{
                        Id:    1,
                        Name:  "Test Product",
                        Images: []string{
                            "test.jpg",
                        },
                    },
                }, nil)
            },
            wantErr: nil,
        },
        {
            name: "Out of stock",
            req: &cart.AddItemRequest{
                UserId:    1,
                ProductId: 1,
                SkuId:    2,
                Quantity: 5,
            },
            mock: func() {
                mockProduct.EXPECT().GetSku(context.Background(), &product.GetSkuRequest{
                    Id: 2,
                }).Return(&product.GetSkuResponse{
                    Sku: &product.Sku{
                        Id:      2,
                        SkuCode: "SKU002",
                        Price:   99.9,
                        Stock:   3,
                    },
                }, nil)
            },
            wantErr: zeroerr.ErrItemOutOfStock,
        },
        {
            name: "Invalid input",
            req: &cart.AddItemRequest{
                UserId:    0,
                ProductId: 1,
                SkuId:    1,
                Quantity: 1,
            },
            mock:    func() {},
            wantErr: zeroerr.ErrInvalidParam,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            tt.mock()

            l := NewAddItemLogic(context.Background(), ctx)
            resp, err := l.AddItem(tt.req)

            if tt.wantErr != nil {
                assert.Error(t, err)
                assert.Equal(t, tt.wantErr, err)
                assert.Nil(t, resp)
            } else {
                assert.NoError(t, err)
                assert.NotNil(t, resp)
                assert.True(t, resp.Success)

                // Verify cart item and stats
                item, err := ctx.CartItemsModel.FindOneByUserIdSkuId(context.Background(), uint64(tt.req.UserId), uint64(tt.req.SkuId))
                if assert.NoError(t, err) {
                    assert.Equal(t, int64(tt.req.Quantity), item.Quantity)
                    assert.Equal(t, "Test Product", item.ProductName)
                    assert.Equal(t, "SKU001", item.SkuName)
                    assert.Equal(t, 99.9, item.Price)
                }

                stats, err := ctx.CartStatsModel.FindOne(context.Background(), uint64(tt.req.UserId))
                if assert.NoError(t, err) {
                    assert.Equal(t, int64(tt.req.Quantity), stats.TotalQuantity)
                    assert.Equal(t, 99.9*float64(tt.req.Quantity), stats.TotalAmount)
                }
            }
        })
    }

    // Cleanup
    _ = ctx.CartItemsModel.DeleteByUserId(context.Background(), 1)
    _ = ctx.CartStatsModel.Delete(context.Background(), 1)
}