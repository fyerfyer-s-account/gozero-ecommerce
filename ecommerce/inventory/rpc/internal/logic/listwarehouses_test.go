package logic

import (
    "context"
    "database/sql"
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

type ListWarehousesTestSuite struct {
    suite.Suite
    ctx    *svc.ServiceContext
    logic  *ListWarehousesLogic
}

func TestListWarehousesSuite(t *testing.T) {
    suite.Run(t, new(ListWarehousesTestSuite))
}

func (s *ListWarehousesTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/inventory.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewListWarehousesLogic(context.Background(), s.ctx)
}

func (s *ListWarehousesTestSuite) SetupTest() {
    s.cleanData()
    s.setupTestData()
}

func (s *ListWarehousesTestSuite) TearDownTest() {
    s.cleanData()
}

func (s *ListWarehousesTestSuite) cleanData() {
    _ = s.ctx.WarehousesModel.Delete(context.Background(), 1)
    _ = s.ctx.WarehousesModel.Delete(context.Background(), 2)
    _ = s.ctx.WarehousesModel.Delete(context.Background(), 3)
}

func (s *ListWarehousesTestSuite) setupTestData() {
    testWarehouses := []*model.Warehouses{
        {
            Name:    "Warehouse 1",
            Address: "Address 1",
            Contact: sql.NullString{String: "Contact 1", Valid: true},
            Phone:   sql.NullString{String: "1234567890", Valid: true},
            Status:  1,
        },
        {
            Name:    "Warehouse 2",
            Address: "Address 2",
            Contact: sql.NullString{String: "Contact 2", Valid: true},
            Phone:   sql.NullString{String: "0987654321", Valid: true},
            Status:  1,
        },
        {
            Name:    "Warehouse 3",
            Address: "Address 3",
            Status:  2,
        },
    }

    for _, w := range testWarehouses {
        _, err := s.ctx.WarehousesModel.Insert(context.Background(), w)
        assert.NoError(s.T(), err)
    }
}

func (s *ListWarehousesTestSuite) TestListWarehouses() {
    tests := []struct {
        name    string
        req     *inventory.ListWarehousesRequest
        check   func(*inventory.ListWarehousesResponse)
    }{
        {
            name: "list all warehouses",
            req: &inventory.ListWarehousesRequest{
                Page:     1,
                PageSize: 10,
            },
            check: func(resp *inventory.ListWarehousesResponse) {
                assert.Equal(s.T(), int64(3), resp.Total)
                assert.Len(s.T(), resp.Warehouses, 3)
            },
        },
        {
            name: "pagination test",
            req: &inventory.ListWarehousesRequest{
                Page:     1,
                PageSize: 2,
            },
            check: func(resp *inventory.ListWarehousesResponse) {
                assert.Equal(s.T(), int64(3), resp.Total)
                assert.Len(s.T(), resp.Warehouses, 2)
            },
        },
        {
            name: "empty page test",
            req: &inventory.ListWarehousesRequest{
                Page:     100,
                PageSize: 10,
            },
            check: func(resp *inventory.ListWarehousesResponse) {
                assert.Equal(s.T(), int64(3), resp.Total)
                assert.Len(s.T(), resp.Warehouses, 0)
            },
        },
        {
            name: "default parameters test",
            req: &inventory.ListWarehousesRequest{},
            check: func(resp *inventory.ListWarehousesResponse) {
                assert.Equal(s.T(), int64(3), resp.Total)
                assert.Len(s.T(), resp.Warehouses, 3)
            },
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            resp, err := s.logic.ListWarehouses(tt.req)
            assert.NoError(t, err)
            tt.check(resp)
        })
    }
}