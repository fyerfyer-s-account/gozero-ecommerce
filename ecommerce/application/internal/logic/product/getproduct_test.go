package product

import (
	"context"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestGetProductLogic_GetProduct(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.GetProductReq
		mock    func(mockProduct *ProductService)
		want    *types.Product
		wantErr error
	}{
		{
			name: "Success",
			req: &types.GetProductReq{
				Id: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().GetProduct(
					mock.Anything,
					&product.GetProductRequest{Id: 1},
				).Return(&product.GetProductResponse{
					Product: &product.Product{
						Id:          1,
						Name:        "Test Product",
						Description: "Test Description",
						CategoryId:  1,
						Brand:       "Test Brand",
						Images:      []string{"image1.jpg"},
						Price:       99.99,
						Sales:       10,
						Status:      1,
						CreatedAt:   1609459200, // 2021-01-01
					},
					Skus: []*product.Sku{
						{
							Id:        1,
							ProductId: 1,
							SkuCode:   "SKU001",
							Price:     99.99,
							Stock:     100,
						},
					},
				}, nil)
			},
			want: &types.Product{
				Id:          1,
				Name:        "Test Product",
				Description: "Test Description",
				CategoryId:  1,
				Brand:       "Test Brand",
				Images:      []string{"image1.jpg"},
				Price:       99.99,
				Stock:       100,
				Sales:       10,
				Status:      1,
				CreatedAt:   1609459200,
			},
			wantErr: nil,
		},
		{
			name: "Invalid ID",
			req: &types.GetProductReq{
				Id: 0,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().GetProduct(
					mock.Anything,
					&product.GetProductRequest{Id: 0},
				).Return(nil, zeroerr.ErrInvalidParam)
			},
			want:    nil,
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Product not found",
			req: &types.GetProductReq{
				Id: 999,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().GetProduct(
					mock.Anything,
					&product.GetProductRequest{Id: 999},
				).Return(nil, zeroerr.ErrProductNotFound)
			},
			want:    nil,
			wantErr: zeroerr.ErrProductNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			logic := NewGetProductLogic(context.Background(), svcCtx)
			got, err := logic.GetProduct(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
				assert.Nil(t, got)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
