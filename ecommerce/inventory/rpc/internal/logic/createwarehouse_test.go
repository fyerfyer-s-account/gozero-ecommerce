package logic

import (
    "context"
    "flag"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/suite"
    "github.com/zeromicro/go-zero/core/conf"
)

type WarehouseTestSuite struct {
    suite.Suite
    ctx   *svc.ServiceContext
    logic *CreateWarehouseLogic
}

func TestWarehouseSuite(t *testing.T) {
    suite.Run(t, new(WarehouseTestSuite))
}

func (s *WarehouseTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/inventory.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewCreateWarehouseLogic(context.Background(), s.ctx)
}

func (s *WarehouseTestSuite) TearDownTest() {
    // Clean up test data after each test
    err := s.ctx.WarehousesModel.Delete(context.Background(), 1)
    assert.NoError(s.T(), err)
}

func (s *WarehouseTestSuite) TestCreateWarehouse() {
    tests := []struct {
        name    string
        req     *inventory.CreateWarehouseRequest
        wantErr bool
        check   func(*testing.T, *inventory.CreateWarehouseResponse)
    }{
        {
            name: "normal case",
            req: &inventory.CreateWarehouseRequest{
                Name:    "Test Warehouse",
                Address: "Test Address",
                Contact: "Test Contact",
                Phone:   "12345678901",
            },
            wantErr: false,
            check: func(t *testing.T, resp *inventory.CreateWarehouseResponse) {
                assert.Greater(t, resp.Id, int64(0))
                
                // Verify warehouse exists
                warehouse, err := s.ctx.WarehousesModel.FindOne(context.Background(), uint64(resp.Id))
                assert.NoError(t, err)
                assert.Equal(t, "Test Warehouse", warehouse.Name)
                assert.Equal(t, "Test Address", warehouse.Address)
                assert.Equal(t, "Test Contact", warehouse.Contact.String)
                assert.Equal(t, "12345678901", warehouse.Phone.String)
                assert.Equal(t, int64(1), warehouse.Status)
            },
        },
        {
            name: "empty name",
            req: &inventory.CreateWarehouseRequest{
                Name:    "",
                Address: "Test Address",
            },
            wantErr: true,
        },
        {
            name: "empty address",
            req: &inventory.CreateWarehouseRequest{
                Name:    "Test Warehouse",
                Address: "",
            },
            wantErr: true,
        },
        {
            name: "duplicate name",
            req: &inventory.CreateWarehouseRequest{
                Name:    "Test Warehouse",
                Address: "Another Address",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            resp, err := s.logic.CreateWarehouse(tt.req)
            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
                if tt.check != nil {
                    tt.check(t, resp)
                }
            }
        })
    }
}