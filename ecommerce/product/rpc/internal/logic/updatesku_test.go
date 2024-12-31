package logic

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/suite"
	"github.com/zeromicro/go-zero/core/conf"
)

type UpdateSkuSuite struct {
	suite.Suite
	ctx         context.Context
	serviceCtx  *svc.ServiceContext
	testProduct *model.Products
	testSku     *model.Skus
}

func (suite *UpdateSkuSuite) SetupSuite() {
	// Load config and initialize service context
	var c config.Config
	conf.MustLoad("../../etc/product.yaml", &c)
	suite.serviceCtx = svc.NewServiceContext(c)
	suite.ctx = context.Background()

	// Insert test product
	testProduct := &model.Products{
		Name:   "Test Product",
		Price:  99.99,
		Status: 1,
	}
	result, err := suite.serviceCtx.ProductsModel.Insert(suite.ctx, testProduct)
	suite.Require().NoError(err)
	productId, err := result.LastInsertId()
	suite.Require().NoError(err)
	suite.testProduct = testProduct
	suite.testProduct.Id = uint64(productId)

	// Insert test SKU
	defaultAttrs := []map[string]string{
		{"color": "red", "size": "M"},
	}
	attrsJSON, err := json.Marshal(defaultAttrs)
	suite.Require().NoError(err)

	testSku := &model.Skus{
		ProductId:  suite.testProduct.Id,
		SkuCode:    "TEST-SKU-001",
		Price:      99.99,
		Stock:      100,
		Attributes: string(attrsJSON),
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	result, err = suite.serviceCtx.SkusModel.Insert(suite.ctx, testSku)
	suite.Require().NoError(err)
	skuId, err := result.LastInsertId()
	suite.Require().NoError(err)
	suite.testSku = testSku
	suite.testSku.Id = uint64(skuId)
}

func (suite *UpdateSkuSuite) TearDownSuite() {
	// Cleanup test data
	_ = suite.serviceCtx.SkusModel.Delete(suite.ctx, suite.testSku.Id)
	_ = suite.serviceCtx.ProductsModel.Delete(suite.ctx, suite.testProduct.Id)
}

func (suite *UpdateSkuSuite) TestUpdateSku() {
	tests := []struct {
		name    string
		req     *product.UpdateSkuRequest
		wantErr error
	}{
		{
			name: "Valid update",
			req: &product.UpdateSkuRequest{
				Id:    int64(suite.testSku.Id),
				Price: 199.99,
				Stock: 50,
			},
			wantErr: nil,
		},
		{
			name: "Invalid price",
			req: &product.UpdateSkuRequest{
				Id:    int64(suite.testSku.Id),
				Price: 0,
				Stock: 100,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Invalid stock",
			req: &product.UpdateSkuRequest{
				Id:    int64(suite.testSku.Id),
				Price: 99.99,
				Stock: -1,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Non-existent SKU",
			req: &product.UpdateSkuRequest{
				Id:    99999,
				Price: 99.99,
				Stock: 100,
			},
			wantErr: zeroerr.ErrSkuNotFound,
		},
	}

	for _, tt := range tests {
		suite.Run(tt.name, func() {
			l := NewUpdateSkuLogic(suite.ctx, suite.serviceCtx)
			resp, err := l.UpdateSku(tt.req)

			if tt.wantErr != nil {
				suite.Error(err)
				suite.Equal(tt.wantErr, err)
				suite.Nil(resp)
			} else {
				suite.NoError(err)
				suite.NotNil(resp)
				suite.True(resp.Success)

				// Verify changes
				updated, err := suite.serviceCtx.SkusModel.FindOne(suite.ctx, suite.testSku.Id)
				suite.NoError(err)
				suite.Equal(tt.req.Price, updated.Price)
				suite.Equal(tt.req.Stock, updated.Stock)
			}
		})
	}
}

func TestUpdateSkuSuite(t *testing.T) {
	suite.Run(t, new(UpdateSkuSuite))
}
