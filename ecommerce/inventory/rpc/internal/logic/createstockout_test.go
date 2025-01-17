package logic

import (
    "context"
    "flag"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type StockOutTestSuite struct {
    suite.Suite
    ctx    *svc.ServiceContext
    logic  *CreateStockOutLogic
    stocks []*model.Stocks
}

func TestStockOutSuite(t *testing.T) {
    suite.Run(t, new(StockOutTestSuite))
}

func (s *StockOutTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/inventory.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewCreateStockOutLogic(context.Background(), s.ctx)
}

func (s *StockOutTestSuite) SetupTest() {
    // Clear existing test data first
    err := s.ctx.StocksModel.Trans(context.Background(), func(ctx context.Context, session sqlx.Session) error {
        _, err := session.Exec("DELETE FROM stocks WHERE warehouse_id = 1")
        if err != nil {
            return err
        }
        _, err = session.Exec("DELETE FROM stock_records WHERE warehouse_id = 1")
        return err
    })
    assert.NoError(s.T(), err)

    // Insert test data
    s.stocks = []*model.Stocks{
        {
            SkuId:         1001,
            WarehouseId:   1,
            Available:     200,
            Locked:        0,
            Total:         200,
            AlertQuantity: 10,
        },
        {
            SkuId:         1002,
            WarehouseId:   1,
            Available:     300,
            Locked:        0,
            Total:         300,
            AlertQuantity: 20,
        },
    }

    for _, stock := range s.stocks {
        _, err := s.ctx.StocksModel.Insert(context.Background(), stock)
        assert.NoError(s.T(), err)
    }
}

func (s *StockOutTestSuite) TearDownTest() {
    // Clean up test data
    err := s.ctx.StocksModel.Trans(context.Background(), func(ctx context.Context, session sqlx.Session) error {
        _, err := session.Exec("DELETE FROM stocks WHERE warehouse_id = 1")
        if err != nil {
            return err
        }
        _, err = session.Exec("DELETE FROM stock_records WHERE warehouse_id = 1")
        return err
    })
    assert.NoError(s.T(), err)
}

func (s *StockOutTestSuite) TestCreateStockOut() {
    tests := []struct {
        name    string
        req     *inventory.CreateStockOutRequest
        wantErr bool
    }{
        {
            name: "normal case",
            req: &inventory.CreateStockOutRequest{
                WarehouseId: 1,
                Items: []*inventory.StockOutItem{
                    {
                        SkuId:    1001,
                        Quantity: 50,
                    },
                    {
                        SkuId:    1002,
                        Quantity: 100,
                    },
                },
                OrderNo: "TEST123",
                Remark:  "Test stock out",
            },
            wantErr: false,
        },
        {
            name: "insufficient stock",
            req: &inventory.CreateStockOutRequest{
                WarehouseId: 1,
                Items: []*inventory.StockOutItem{
                    {
                        SkuId:    1001,
                        Quantity: 1000, // More than available
                    },
                },
            },
            wantErr: true,
        },
        {
            name: "invalid warehouse",
            req: &inventory.CreateStockOutRequest{
                WarehouseId: 0,
                Items: []*inventory.StockOutItem{
                    {
                        SkuId:    1001,
                        Quantity: 50,
                    },
                },
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            resp, err := s.logic.CreateStockOut(tt.req)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                assert.True(t, resp.Success)

                // Verify stock quantities
                for _, item := range tt.req.Items {
                    stock, err := s.ctx.StocksModel.FindOneBySkuIdWarehouseId(context.Background(), 
                        uint64(item.SkuId), uint64(tt.req.WarehouseId))
                    assert.NoError(t, err)
                    assert.GreaterOrEqual(t, stock.Available, int64(0))
                }
            }
        })
    }
}