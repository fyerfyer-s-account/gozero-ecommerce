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
)

type UnlockStockTestSuite struct {
    suite.Suite
    ctx    *svc.ServiceContext
    logic  *UnlockStockLogic
}

func TestUnlockStockSuite(t *testing.T) {
    suite.Run(t, new(UnlockStockTestSuite))
}

func (s *UnlockStockTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/inventory.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewUnlockStockLogic(context.Background(), s.ctx)
}

func (s *UnlockStockTestSuite) SetupTest() {
    s.cleanData()
    s.setupTestData()
}

func (s *UnlockStockTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *UnlockStockTestSuite) cleanData() {
    _ = s.ctx.StocksModel.Delete(context.Background(), 1)
    _ = s.ctx.StockLocksModel.DeleteByOrderNo(context.Background(), "TEST123")
}

func (s *UnlockStockTestSuite) setupTestData() {
    // Create stock
    stock := &model.Stocks{
        SkuId:       1001,
        WarehouseId: 1,
        Available:   150,
        Locked:      50,
        Total:       200,
    }
    _, err := s.ctx.StocksModel.Insert(context.Background(), stock)
    assert.NoError(s.T(), err)

    // Create lock record
    lock := &model.StockLocks{
        OrderNo:     "TEST123",
        SkuId:       1001,
        WarehouseId: 1,
        Quantity:    50,
        Status:      1,
    }
    _, err = s.ctx.StockLocksModel.Insert(context.Background(), lock)
    assert.NoError(s.T(), err)
}

func (s *UnlockStockTestSuite) TestUnlockStock() {
    tests := []struct {
        name    string
        req     *inventory.UnlockStockRequest
        wantErr bool
        check   func()
    }{
        {
            name: "normal case",
            req: &inventory.UnlockStockRequest{
                OrderNo: "TEST123",
            },
            wantErr: false,
            check: func() {
                // Verify stock
                stock, err := s.ctx.StocksModel.FindOneBySkuIdWarehouseId(context.Background(), 1001, 1)
                assert.NoError(s.T(), err)
                assert.Equal(s.T(), int64(200), stock.Available)
                assert.Equal(s.T(), int64(0), stock.Locked)

                // Verify lock record deleted
                locks, err := s.ctx.StockLocksModel.FindByOrderNo(context.Background(), "TEST123")
                assert.NoError(s.T(), err)
                assert.Empty(s.T(), locks)
            },
        },
        {
            name: "empty order number",
            req: &inventory.UnlockStockRequest{
                OrderNo: "",
            },
            wantErr: true,
        },
        {
            name: "non-existent locks",
            req: &inventory.UnlockStockRequest{
                OrderNo: "NONEXIST",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            resp, err := s.logic.UnlockStock(tt.req)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.True(t, resp.Success)
            if tt.check != nil {
                tt.check()
            }
        })
    }
}