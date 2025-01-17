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

type UpdateWarehouseTestSuite struct {
    suite.Suite
    ctx    *svc.ServiceContext
    logic  *UpdateWarehouseLogic
    testID uint64
}

func TestUpdateWarehouseSuite(t *testing.T) {
    suite.Run(t, new(UpdateWarehouseTestSuite))
}

func (s *UpdateWarehouseTestSuite) SetupSuite() {
    configFile := flag.String("f", "../../etc/inventory.yaml", "the config file")
    var c config.Config
    conf.MustLoad(*configFile, &c)
    s.ctx = svc.NewServiceContext(c)
    s.logic = NewUpdateWarehouseLogic(context.Background(), s.ctx)
}

func (s *UpdateWarehouseTestSuite) SetupTest() {
    // Create test warehouse
    warehouse := &model.Warehouses{
        Name:      "Test Warehouse",
        Address:   "Test Address",
        Contact:   sql.NullString{String: "Test Contact", Valid: true},
        Phone:     sql.NullString{String: "12345678901", Valid: true},
        Status:    1,
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    result, err := s.ctx.WarehousesModel.Insert(context.Background(), warehouse)
    assert.NoError(s.T(), err)
    id, err := result.LastInsertId()
    assert.NoError(s.T(), err)
    s.testID = uint64(id)
}

func (s *UpdateWarehouseTestSuite) TearDownTest() {
    // Clean up test data
    _ = s.ctx.WarehousesModel.Delete(context.Background(), s.testID)
}

func (s *UpdateWarehouseTestSuite) TestUpdateWarehouse() {
    tests := []struct {
        name    string
        req     *inventory.UpdateWarehouseRequest
        wantErr bool
        check   func(*model.Warehouses)
    }{
        {
            name: "update all fields",
            req: &inventory.UpdateWarehouseRequest{
                Id:      int64(s.testID),
                Name:    "Updated Warehouse",
                Address: "Updated Address",
                Contact: "Updated Contact",
                Phone:   "98765432101",
                Status:  2,
            },
            wantErr: false,
            check: func(w *model.Warehouses) {
                assert.Equal(s.T(), "Updated Warehouse", w.Name)
                assert.Equal(s.T(), "Updated Address", w.Address)
                assert.Equal(s.T(), "Updated Contact", w.Contact.String)
                assert.Equal(s.T(), "98765432101", w.Phone.String)
                assert.Equal(s.T(), int64(2), w.Status)
            },
        },
        {
            name: "update partial fields",
            req: &inventory.UpdateWarehouseRequest{
                Id:      int64(s.testID),
                Name:    "New Name",
                Address: "New Address",
            },
            wantErr: false,
            check: func(w *model.Warehouses) {
                assert.Equal(s.T(), "New Name", w.Name)
                assert.Equal(s.T(), "New Address", w.Address)
            },
        },
        {
            name: "invalid id",
            req: &inventory.UpdateWarehouseRequest{
                Id: 0,
            },
            wantErr: true,
        },
        {
            name: "non-existent warehouse",
            req: &inventory.UpdateWarehouseRequest{
                Id:   999999,
                Name: "Non-existent",
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        s.T().Run(tt.name, func(t *testing.T) {
            resp, err := s.logic.UpdateWarehouse(tt.req)
            if tt.wantErr {
                assert.Error(t, err)
                return
            }
            assert.NoError(t, err)
            assert.True(t, resp.Success)

            if tt.check != nil {
                warehouse, err := s.ctx.WarehousesModel.FindOne(context.Background(), s.testID)
                assert.NoError(t, err)
                tt.check(warehouse)
            }
        })
    }
}