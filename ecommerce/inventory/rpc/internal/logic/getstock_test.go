package logic

import (
    "context"
    "flag"
    "testing"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/zeromicro/go-zero/core/conf"
)

type GetStockTestSuite struct {
    suite.Suite
    ctx    *svc.ServiceContext
    logic  *GetStockLogic
}

func TestGetStockSuite(t *testing.T) {
    suite.Run(t, new(GetStockTestSuite))
}

func (s *GetStockTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/inventory.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewGetStockLogic(context.Background(), s.ctx)
}

func (s *GetStockTestSuite) SetupTest() {
    // Clean existing data
    s.cleanData()

    // Insert test data
    stock := &model.Stocks{
        SkuId:         1001,
        WarehouseId:   1,
        Available:     100,
        Locked:        20,
        Total:         120,
        AlertQuantity: 10,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }

    _, err := s.ctx.StocksModel.Insert(context.Background(), stock)
    assert.NoError(s.T(), err)
}

func (s *GetStockTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *GetStockTestSuite) cleanData() {
    _ = s.ctx.StocksModel.Delete(context.Background(), 1)
}

func (s *GetStockTestSuite) TestGetStock() {
    tests := []struct {
        name    string
        req     *inventory.GetStockRequest
        wantErr bool
        check   func(*inventory.GetStockResponse)
    }{
        {
            name: "normal case",
            req: &inventory.GetStockRequest{
                SkuId:       1001,
                WarehouseId: 1,
            },
            wantErr: false,
            check: func(resp *inventory.GetStockResponse) {
                assert.NotNil(s.T(), resp.Stock)
                assert.Equal(s.T(), int64(1001), resp.Stock.SkuId)
                assert.Equal(s.T(), int64(1), resp.Stock.WarehouseId)
                assert.Equal(s.T(), int32(100), resp.Stock.Available)
                assert.Equal(s.T(), int32(20), resp.Stock.Locked)
                assert.Equal(s.T(), int32(120), resp.Stock.Total)
            },
        },
        {
            name: "not found",
            req: &inventory.GetStockRequest{
                SkuId:       1002,
                WarehouseId: 1,
            },
            wantErr: true,
        },
        {
            name: "invalid sku id",
            req: &inventory.GetStockRequest{
                SkuId:       0,
                WarehouseId: 1,
            },
            wantErr: true,
        },
        {
            name: "invalid warehouse id",
            req: &inventory.GetStockRequest{
                SkuId:       1001,
                WarehouseId: 0,
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            resp, err := s.logic.GetStock(tt.req)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            if tt.check != nil {
                tt.check(resp)
            }
        })
    }
}