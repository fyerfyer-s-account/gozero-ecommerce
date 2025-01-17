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

type UpdateStockTestSuite struct {
	suite.Suite
	ctx   *svc.ServiceContext
	logic *UpdateStockLogic
}

func TestUpdateStockSuite(t *testing.T) {
	suite.Run(t, new(UpdateStockTestSuite))
}

func (s *UpdateStockTestSuite) SetupSuite() {
	configFile := flag.String("f", "../../etc/inventory.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
	s.logic = NewUpdateStockLogic(context.Background(), s.ctx)
}

func (s *UpdateStockTestSuite) SetupTest() {
	s.cleanData()
	s.setupTestData()
}

func (s *UpdateStockTestSuite) TearDownTest() {
	s.cleanData()
}

func (s *UpdateStockTestSuite) cleanData() {
	ctx := context.Background()
	_ = s.ctx.StocksModel.Trans(ctx, func(ctx context.Context, session sqlx.Session) error {
		_, _ = session.Exec("DELETE FROM stocks WHERE warehouse_id = 1")
		_, _ = session.Exec("DELETE FROM stock_records WHERE warehouse_id = 1")
		return nil
	})
}

func (s *UpdateStockTestSuite) setupTestData() {
	stock := &model.Stocks{
		SkuId:       1001,
		WarehouseId: 1,
		Available:   100,
		Locked:      0,
		Total:       100,
	}
	_, err := s.ctx.StocksModel.Insert(context.Background(), stock)
	assert.NoError(s.T(), err)
}

func (s *UpdateStockTestSuite) TestUpdateStock() {
	tests := []struct {
		name    string
		req     *inventory.UpdateStockRequest
		wantErr bool
		check   func()
	}{
		{
			name: "increase stock",
			req: &inventory.UpdateStockRequest{
				SkuId:       1001,
				WarehouseId: 1,
				Quantity:    50,
				Remark:      "test increase",
			},
			wantErr: false,
			check: func() {
				stock, err := s.ctx.StocksModel.FindOneBySkuIdWarehouseId(context.Background(), 1001, 1)
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), int64(150), stock.Available)
				assert.Equal(s.T(), int64(150), stock.Total)
			},
		},
		{
			name: "decrease stock",
			req: &inventory.UpdateStockRequest{
				SkuId:       1001,
				WarehouseId: 1,
				Quantity:    -30,
				Remark:      "test decrease",
			},
			wantErr: false,
			check: func() {
				stock, err := s.ctx.StocksModel.FindOneBySkuIdWarehouseId(context.Background(), 1001, 1)
				assert.NoError(s.T(), err)
				assert.Equal(s.T(), int64(70), stock.Available)
				assert.Equal(s.T(), int64(70), stock.Total)
			},
		},
		{
			name: "insufficient stock",
			req: &inventory.UpdateStockRequest{
				SkuId:       1001,
				WarehouseId: 1,
				Quantity:    -200,
				Remark:      "test insufficient",
			},
			wantErr: true,
		},
		{
			name: "invalid sku",
			req: &inventory.UpdateStockRequest{
				SkuId:       0,
				WarehouseId: 1,
				Quantity:    50,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			// Reset test data for each case
			s.cleanData()
			s.setupTestData()

			resp, err := s.logic.UpdateStock(tt.req)
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
