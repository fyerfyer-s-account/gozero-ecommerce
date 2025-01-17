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
	"github.com/zeromicro/go-zero/core/conf"
)

func TestBatchGetStockLogic(t *testing.T) {
	configFile := flag.String("f", "../../etc/inventory.yaml", "the config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	
	// Create service context
	ctx := svc.NewServiceContext(c)
	
	// Setup test data
	setupTestData(t, ctx)
	
	logic := NewBatchGetStockLogic(context.Background(), ctx)

	tests := []struct {
		name    string
		req     *inventory.BatchGetStockRequest
		wantErr bool
		check   func(*testing.T, *inventory.BatchGetStockResponse)
	}{
		{
			name: "normal case",
			req: &inventory.BatchGetStockRequest{
				SkuIds:      []int64{1001, 1002},
				WarehouseId: 1,
			},
			wantErr: false,
			check: func(t *testing.T, resp *inventory.BatchGetStockResponse) {
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Stocks)
				// We only check if response is properly formatted
				// since test data might vary
				for _, stock := range resp.Stocks {
					assert.NotNil(t, stock)
					assert.Greater(t, stock.Id, int64(0))
					assert.Greater(t, stock.SkuId, int64(0))
					assert.Greater(t, stock.WarehouseId, int64(0))
					assert.GreaterOrEqual(t, stock.Total, stock.Available+stock.Locked)
				}
			},
		},
		{
			name: "empty sku list",
			req: &inventory.BatchGetStockRequest{
				SkuIds:      []int64{},
				WarehouseId: 1,
			},
			wantErr: false,
			check: func(t *testing.T, resp *inventory.BatchGetStockResponse) {
				assert.NotNil(t, resp)
				assert.NotNil(t, resp.Stocks)
				assert.Equal(t, 0, len(resp.Stocks))
			},
		},
		{
			name: "invalid sku id",
			req: &inventory.BatchGetStockRequest{
				SkuIds:      []int64{-1},
				WarehouseId: 1,
			},
			wantErr: true,
			check: func(t *testing.T, resp *inventory.BatchGetStockResponse) {
				assert.Nil(t, resp)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp, err := logic.BatchGetStock(tt.req)
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

func setupTestData(t *testing.T, ctx *svc.ServiceContext) {
	// Insert test stock data
	stocks := []*model.Stocks{
		{
			SkuId:       1001,
			WarehouseId: 1,
			Available:   100,
			Locked:      0,
			Total:       100,
			AlertQuantity: 10,
		},
		{
			SkuId:       1002,
			WarehouseId: 1,
			Available:   200,
			Locked:      50,
			Total:       250,
			AlertQuantity: 20,
		},
	}

	for _, stock := range stocks {
		_, err := ctx.StocksModel.Insert(context.Background(), stock)
		assert.NoError(t, err)
	}
}
