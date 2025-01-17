package logic

import (
    "context"
    "database/sql"
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

type ListStockRecordsTestSuite struct {
    suite.Suite
    ctx    *svc.ServiceContext
    logic  *ListStockRecordsLogic
}

func TestListStockRecordsSuite(t *testing.T) {
    suite.Run(t, new(ListStockRecordsTestSuite))
}

func (s *ListStockRecordsTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/inventory.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewListStockRecordsLogic(context.Background(), s.ctx)
}

func (s *ListStockRecordsTestSuite) SetupTest() {
    s.cleanData()
    s.setupTestData()
}

func (s *ListStockRecordsTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *ListStockRecordsTestSuite) cleanData() {
    _ = s.ctx.StockRecordsModel.Delete(context.Background(), 1)
}

func (s *ListStockRecordsTestSuite) setupTestData() {
    // Insert test records
    testRecords := []*model.StockRecords{
        {
            SkuId:       1001,
            WarehouseId: 1,
            Type:        1, // Stock in
            Quantity:    100,
            OrderNo:     sql.NullString{String: "ORDER001", Valid: true},
            Remark:      sql.NullString{String: "Test stock in", Valid: true},
            CreatedAt:   time.Now(),
        },
        {
            SkuId:       1001,
            WarehouseId: 1,
            Type:        2, // Stock out
            Quantity:    50,
            OrderNo:     sql.NullString{String: "ORDER002", Valid: true},
            Remark:      sql.NullString{String: "Test stock out", Valid: true},
            CreatedAt:   time.Now(),
        },
        {
            SkuId:       1002,
            WarehouseId: 2,
            Type:        1,
            Quantity:    200,
            OrderNo:     sql.NullString{String: "ORDER003", Valid: true},
            Remark:      sql.NullString{String: "Test stock in 2", Valid: true},
            CreatedAt:   time.Now(),
        },
    }

    for _, record := range testRecords {
        _, err := s.ctx.StockRecordsModel.Insert(context.Background(), record)
        assert.NoError(s.T(), err)
    }
}

func (s *ListStockRecordsTestSuite) TestListStockRecords() {
    tests := []struct {
        name    string
        req     *inventory.ListStockRecordsRequest
        check   func(*inventory.ListStockRecordsResponse)
    }{
        {
            name: "list all records",
            req: &inventory.ListStockRecordsRequest{
                Page:     1,
                PageSize: 10,
            },
            check: func(resp *inventory.ListStockRecordsResponse) {
                assert.Equal(s.T(), int64(3), resp.Total)
                assert.Len(s.T(), resp.Records, 3)
            },
        },
        {
            name: "filter by sku_id",
            req: &inventory.ListStockRecordsRequest{
                SkuId:    1001,
                Page:     1,
                PageSize: 10,
            },
            check: func(resp *inventory.ListStockRecordsResponse) {
                assert.Equal(s.T(), int64(2), resp.Total)
                assert.Len(s.T(), resp.Records, 2)
                for _, record := range resp.Records {
                    assert.Equal(s.T(), int64(1001), record.SkuId)
                }
            },
        },
        {
            name: "filter by warehouse_id",
            req: &inventory.ListStockRecordsRequest{
                WarehouseId: 2,
                Page:        1,
                PageSize:    10,
            },
            check: func(resp *inventory.ListStockRecordsResponse) {
                assert.Equal(s.T(), int64(1), resp.Total)
                assert.Len(s.T(), resp.Records, 1)
                assert.Equal(s.T(), int64(2), resp.Records[0].WarehouseId)
            },
        },
        {
            name: "filter by type",
            req: &inventory.ListStockRecordsRequest{
                Type:     2,
                Page:     1,
                PageSize: 10,
            },
            check: func(resp *inventory.ListStockRecordsResponse) {
                assert.Equal(s.T(), int64(1), resp.Total)
                assert.Len(s.T(), resp.Records, 1)
                assert.Equal(s.T(), int32(2), resp.Records[0].Type)
            },
        },
        {
            name: "pagination",
            req: &inventory.ListStockRecordsRequest{
                Page:     1,
                PageSize: 2,
            },
            check: func(resp *inventory.ListStockRecordsResponse) {
                assert.Equal(s.T(), int64(3), resp.Total)
                assert.Len(s.T(), resp.Records, 2)
            },
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            resp, err := s.logic.ListStockRecords(tt.req)
            assert.NoError(t, err)
            tt.check(resp)
        })
    }
}