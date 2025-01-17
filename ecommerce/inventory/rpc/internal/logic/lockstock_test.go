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

type LockStockTestSuite struct {
    suite.Suite
    ctx    *svc.ServiceContext
    logic  *LockStockLogic
}

func TestLockStockSuite(t *testing.T) {
    suite.Run(t, new(LockStockTestSuite))
}

func (s *LockStockTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/inventory.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewLockStockLogic(context.Background(), s.ctx)
}

func (s *LockStockTestSuite) SetupTest() {
    s.cleanData()
    s.setupTestData()
}

func (s *LockStockTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *LockStockTestSuite) cleanData() {
    _ = s.ctx.StocksModel.Delete(context.Background(), 1)
    _ = s.ctx.StockLocksModel.DeleteByOrderNo(context.Background(), "TEST123")
}

func (s *LockStockTestSuite) setupTestData() {
    stock := &model.Stocks{
        SkuId:       1001,
        WarehouseId: 1,
        Available:   200,
        Locked:      0,
        Total:       200,
    }
    _, err := s.ctx.StocksModel.Insert(context.Background(), stock)
    assert.NoError(s.T(), err)
}

func (s *LockStockTestSuite) TestLockStock() {
    tests := []struct {
        name    string
        req     *inventory.LockStockRequest
        check   func(*inventory.LockStockResponse)
    }{
        {
            name: "normal case",
            req: &inventory.LockStockRequest{
                OrderNo: "TEST123",
                Items: []*inventory.LockItem{
                    {
                        SkuId:       1001,
                        WarehouseId: 1,
                        Quantity:    50,
                    },
                },
            },
            check: func(resp *inventory.LockStockResponse) {
                assert.True(s.T(), resp.Success)
                assert.Empty(s.T(), resp.FailedItems)

                // Verify stock
                stock, err := s.ctx.StocksModel.FindOneBySkuIdWarehouseId(context.Background(), 1001, 1)
                assert.NoError(s.T(), err)
                assert.Equal(s.T(), int64(150), stock.Available)
                assert.Equal(s.T(), int64(50), stock.Locked)

                // Verify lock record
                locks, err := s.ctx.StockLocksModel.FindByOrderNo(context.Background(), "TEST123")
                assert.NoError(s.T(), err)
                assert.Len(s.T(), locks, 1)
                assert.Equal(s.T(), int64(50), locks[0].Quantity)
            },
        },
        {
            name: "insufficient stock",
            req: &inventory.LockStockRequest{
                OrderNo: "TEST123",
                Items: []*inventory.LockItem{
                    {
                        SkuId:       1001,
                        WarehouseId: 1,
                        Quantity:    1000,
                    },
                },
            },
            check: func(resp *inventory.LockStockResponse) {
                assert.False(s.T(), resp.Success)
                assert.Len(s.T(), resp.FailedItems, 1)
                assert.Equal(s.T(), int64(1001), resp.FailedItems[0].SkuId)
            },
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            resp, err := s.logic.LockStock(tt.req)
            assert.NoError(t, err)
            tt.check(resp)
        })
    }
}