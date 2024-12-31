package product

import (
	"context"
	"errors"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestSearchProductsLogic_SearchProducts(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.SearchReq
		mock    func(mockProduct *ProductService)
		want    *types.SearchResp
		wantErr error
	}{
		{
			name: "successful search by keyword",
			req: &types.SearchReq{
				Keyword: "phone",
				Page:    1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ListProducts(
					mock.Anything,
					&product.ListProductsRequest{
						Keyword: "phone",
						Page:    1,
					},
				).Return(&product.ListProductsResponse{
					Total: 1,
					Products: []*product.Product{
						{
							Id:          1,
							Name:        "iPhone",
							Description: "Latest iPhone",
							CategoryId:  1,
							Brand:       "Apple",
							Price:       999.99,
							Sales:       100,
							Status:      1,
							CreatedAt:   1234567890,
						},
					},
				}, nil)
			},
			want: &types.SearchResp{
				List: []types.Product{
					{
						Id:          1,
						Name:        "iPhone",
						Description: "Latest iPhone",
						CategoryId:  1,
						Brand:       "Apple",
						Price:       999.99,
						Sales:       100,
						Status:      1,
						CreatedAt:   1234567890,
					},
				},
				Total:      1,
				Page:       1,
				TotalPages: 1,
			},
			wantErr: nil,
		},
		{
			name: "search by category",
			req: &types.SearchReq{
				CategoryId: 1,
				Page:       1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ListProducts(
					mock.Anything,
					&product.ListProductsRequest{
						CategoryId: 1,
						Page:       1,
					},
				).Return(&product.ListProductsResponse{
					Total:    0,
					Products: []*product.Product{},
				}, nil)
			},
			want: &types.SearchResp{
				List:       []types.Product{},
				Total:      0,
				Page:       1,
				TotalPages: 0,
			},
			wantErr: nil,
		},
		{
			name: "rpc error",
			req: &types.SearchReq{
				Page: 1,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ListProducts(
					mock.Anything,
					&product.ListProductsRequest{
						Page: 1,
					},
				).Return(nil, errors.New("rpc error"))
			},
			want:    nil,
			wantErr: errors.New("rpc error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
				Config: config.Config{
					PageSize: 10,
				},
			}

			logic := NewSearchProductsLogic(context.Background(), svcCtx)
			got, err := logic.SearchProducts(tt.req)

			if tt.wantErr != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.want, got)
			}
		})
	}
}
