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

type DeductStockTestSuite struct {
	suite.Suite
	ctx    *svc.ServiceContext
	logic  *DeductStockLogic
	stocks []*model.Stocks
	locks  []*model.StockLocks
}

func TestDeductStockSuite(t *testing.T) {
	suite.Run(t, new(DeductStockTestSuite))
}

func (s *DeductStockTestSuite) SetupSuite() {
	configFile := flag.String("f", "../../etc/inventory.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	s.ctx = svc.NewServiceContext(c)
	s.logic = NewDeductStockLogic(context.Background(), s.ctx)
}

func (s *DeductStockTestSuite) SetupTest() {
	// Clean existing data
	s.cleanData()

	// Setup test stocks
	s.stocks = []*model.Stocks{
		{
			SkuId:       1001,
			WarehouseId: 1,
			Available:   100,
			Locked:      50,
			Total:       150,
		},
		{
			SkuId:       1002,
			WarehouseId: 1,
			Available:   200,
			Locked:      30,
			Total:       230,
		},
	}

	// Setup test locks
	s.locks = []*model.StockLocks{
		{
			OrderNo:     "TEST123",
			SkuId:       1001,
			WarehouseId: 1,
			Quantity:    30,
			Status:      1,
		},
		{
			OrderNo:     "TEST123",
			SkuId:       1002,
			WarehouseId: 1,
			Quantity:    20,
			Status:      1,
		},
	}

	// Insert test data
	for _, stock := range s.stocks {
		_, err := s.ctx.StocksModel.Insert(context.Background(), stock)
		assert.NoError(s.T(), err)
	}

	for _, lock := range s.locks {
		_, err := s.ctx.StockLocksModel.Insert(context.Background(), lock)
		assert.NoError(s.T(), err)
	}
}

func (s *DeductStockTestSuite) TearDownTest() {
	s.cleanData()
}

func (s *DeductStockTestSuite) cleanData() {
	ctx := context.Background()

	// Clean stocks
	for _, stock := range s.stocks {
		existingStock, err := s.ctx.StocksModel.FindOneBySkuIdWarehouseId(ctx,
			stock.SkuId, stock.WarehouseId)
		if err == nil && existingStock != nil {
			_ = s.ctx.StocksModel.Delete(ctx, existingStock.Id)
		}
	}

	// Clean locks by order number
	if len(s.locks) > 0 {
		_ = s.ctx.StockLocksModel.DeleteByOrderNo(ctx, "TEST123")
	}

	// Clean any remaining records
	for _, stock := range s.stocks {
		records, err := s.ctx.StockRecordsModel.FindBySkuAndWarehouse(ctx,
			stock.SkuId, stock.WarehouseId)
		if err == nil {
			for _, record := range records {
				_ = s.ctx.StockRecordsModel.Delete(ctx, record.Id)
			}
		}
	}
}

func (s *DeductStockTestSuite) TestDeductStock() {
	tests := []struct {
		name    string
		req     *inventory.DeductStockRequest
		wantErr bool
	}{
		{
			name: "normal case",
			req: &inventory.DeductStockRequest{
				OrderNo: "TEST123",
			},
			wantErr: false,
		},
		{
			name: "empty order number",
			req: &inventory.DeductStockRequest{
				OrderNo: "",
			},
			wantErr: true,
		},
		{
			name: "non-existent order",
			req: &inventory.DeductStockRequest{
				OrderNo: "NONEXIST",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			resp, err := s.logic.DeductStock(tt.req)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.True(t, resp.Success)

				 // Fix verification by using correct index
                for i, lock := range s.locks {
                    stock, err := s.ctx.StocksModel.FindOneBySkuIdWarehouseId(context.Background(),
                        lock.SkuId, lock.WarehouseId)
                    assert.NoError(t, err)
                    // Compare with correct initial stock
                    assert.Equal(t, s.stocks[i].Locked-lock.Quantity, stock.Locked)
                }

				// Verify locks are deleted
				locks, err := s.ctx.StockLocksModel.FindByOrderNo(context.Background(), tt.req.OrderNo)
				assert.NoError(t, err)
				assert.Empty(t, locks)
			}
		})
	}
}
