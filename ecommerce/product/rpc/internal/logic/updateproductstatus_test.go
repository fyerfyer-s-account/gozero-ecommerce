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

func TestUpdateProductStatusLogic_UpdateProductStatus(t *testing.T) {
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
		req     *product.UpdateProductStatusRequest
		wantErr error
	}{
		{
			name: "Set status to off shelf",
			req: &product.UpdateProductStatusRequest{
				Id:     productId,
				Status: 2,
			},
			wantErr: nil,
		},
		{
			name: "Set status to on shelf",
			req: &product.UpdateProductStatusRequest{
				Id:     productId,
				Status: 1,
			},
			wantErr: nil,
		},
		{
			name: "Invalid ID",
			req: &product.UpdateProductStatusRequest{
				Id:     0,
				Status: 1,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Invalid status",
			req: &product.UpdateProductStatusRequest{
				Id:     productId,
				Status: 3,
			},
			wantErr: zeroerr.ErrInvalidParam,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := NewUpdateProductStatusLogic(context.Background(), ctx)
			resp, err := l.UpdateProductStatus(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr, err)
				assert.Nil(t, resp)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, resp)
				assert.True(t, resp.Success)

				// Verify status update
				updated, err := ctx.ProductsModel.FindOne(context.Background(), uint64(productId))
				assert.NoError(t, err)
				assert.Equal(t, tt.req.Status, updated.Status)
			}
		})
	}
}
