package product

import (
	"context"
	"errors"
	"testing"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestListCategoriesLogic_ListCategories(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(mockProduct *ProductService)
		want    []types.Category
		wantErr error
	}{
		{
			name: "successful get categories",
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ListCategories(
					mock.Anything,
					&product.ListCategoriesRequest{
						ParentId: 0,
					},
				).Return(&product.ListCategoriesResponse{
					Categories: []*product.Category{
						{
							Id:       1,
							Name:     "Electronics",
							ParentId: 0,
							Level:    1,
							Sort:     1,
							Icon:     "electronics.png",
						},
						{
							Id:       2,
							Name:     "Clothing",
							ParentId: 0,
							Level:    1,
							Sort:     2,
							Icon:     "clothing.png",
						},
					},
				}, nil)
			},
			want: []types.Category{
				{
					Id:       1,
					Name:     "Electronics",
					ParentId: 0,
					Level:    1,
					Sort:     1,
					Icon:     "electronics.png",
				},
				{
					Id:       2,
					Name:     "Clothing",
					ParentId: 0,
					Level:    1,
					Sort:     2,
					Icon:     "clothing.png",
				},
			},
			wantErr: nil,
		},
		{
			name: "empty categories",
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ListCategories(
					mock.Anything,
					&product.ListCategoriesRequest{
						ParentId: 0,
					},
				).Return(&product.ListCategoriesResponse{
					Categories: []*product.Category{},
				}, nil)
			},
			want:    []types.Category{},
			wantErr: nil,
		},
		{
			name: "rpc error",
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().ListCategories(
					mock.Anything,
					&product.ListCategoriesRequest{
						ParentId: 0,
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
			}

			logic := NewListCategoriesLogic(context.Background(), svcCtx)
			got, err := logic.ListCategories()

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
