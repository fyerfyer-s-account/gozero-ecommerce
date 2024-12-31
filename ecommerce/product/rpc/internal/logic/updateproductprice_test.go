package logic

import (
	"context"
	"database/sql"
	"flag"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestUpdateProductPriceLogic_UpdateProductPrice(t *testing.T) {
	configFile := flag.String("f", "../../etc/product.yaml", "config file")
	var c config.Config
	conf.MustLoad(*configFile, &c)
	ctx := svc.NewServiceContext(c)

	// Create test product
	testProduct := &model.Products{
		Name:        "Test Product",
		Description: sql.NullString{String: "Test Description", Valid: true},
		CategoryId:  1,
		Price:       99.99,
		Status:      1,
	}
	result, err := ctx.ProductsModel.Insert(context.Background(), testProduct)
	assert.NoError(t, err)
	productId, err := result.LastInsertId()
	assert.NoError(t, err)

	defer func() {
		_ = ctx.ProductsModel.Delete(context.Background(), uint64(productId))
	}()

	tests := []struct {
		name    string
		req     *product.UpdateProductPriceRequest
		wantErr error
	}{
		{
			name: "Valid price update",
			req: &product.UpdateProductPriceRequest{
				Id:    productId,
				Price: 199.99,
			},
			wantErr: nil,
		},
		{
			name: "Invalid ID",
			req: &product.UpdateProductPriceRequest{
				Id:    0,
				Price: 99.99,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Negative price",
			req: &product.UpdateProductPriceRequest{
				Id:    productId,
				Price: -10.00,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateProductPriceLogic(context.Background(), ctx)
			resp, err := l.UpdateProductPrice(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify price update
				updated, err := ctx.ProductsModel.FindOne(context.Background(), uint64(productId))
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Price, updated.Price)
			}
		})
	}
}
