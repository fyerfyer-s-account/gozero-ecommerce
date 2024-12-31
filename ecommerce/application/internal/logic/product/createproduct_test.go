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

func TestCreateProductLogic_CreateProduct(t *testing.T) {
	tests := []struct {
		name    string
		req     *types.CreateProductReq
		mock    func(mockProduct *ProductService)
		want    *types.CreateProductResp
		wantErr error
	}{
		{
			name: "Success - Basic product",
			req: &types.CreateProductReq{
				Name:        "Test Product",
				Description: "Test Description",
				CategoryId:  1,
				Brand:       "Test Brand",
				Images:      []string{"image1.jpg"},
				Price:       99.99,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().CreateProduct(
					mock.Anything, // Changed from mock.AnythingOfType("*context.backgroundCtx")
					mock.MatchedBy(func(req *product.CreateProductRequest) bool {
						return req.Name == "Test Product" &&
							req.Description == "Test Description" &&
							req.CategoryId == 1 &&
							req.Brand == "Test Brand" &&
							len(req.Images) == 1 &&
							req.Images[0] == "image1.jpg" &&
							req.Price == 99.99
					}),
				).Return(&product.CreateProductResponse{Id: 1}, nil)
			},
			want:    &types.CreateProductResp{Id: 1},
			wantErr: nil,
		},
		{
			name: "Success - With SKU attributes",
			req: &types.CreateProductReq{
				Name:        "Test Product",
				Description: "Test Description",
				CategoryId:  1,
				Brand:       "Test Brand",
				Images:      []string{"image1.jpg"},
				Price:       99.99,
				Attributes: []types.SkuAttributeReq{
					{Key: "Color", Value: "Red"},
					{Key: "Size", Value: "L"},
				},
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().CreateProduct(
					mock.Anything,
					&product.CreateProductRequest{
						Name:        "Test Product",
						Description: "Test Description",
						CategoryId:  1,
						Brand:       "Test Brand",
						Images:      []string{"image1.jpg"},
						Price:       99.99,
						SkuAttributes: []*product.SkuAttribute{
							{Key: "Color", Value: "Red"},
							{Key: "Size", Value: "L"},
						},
					},
				).Return(&product.CreateProductResponse{Id: 1}, nil)
			},
			want:    &types.CreateProductResp{Id: 1},
			wantErr: nil,
		},
		{
			name: "Invalid input",
			req: &types.CreateProductReq{
				Name:       "",
				CategoryId: 1,
				Price:      99.99,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().CreateProduct(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrInvalidParam)
			},
			want:    nil,
			wantErr: zeroerr.ErrInvalidParam,
		},
		{
			name: "Category not found",
			req: &types.CreateProductReq{
				Name:       "Test Product",
				CategoryId: 999,
				Price:      99.99,
			},
			mock: func(mockProduct *ProductService) {
				mockProduct.EXPECT().CreateProduct(
					mock.Anything,
					mock.Anything,
				).Return(nil, zeroerr.ErrCategoryNotFound)
			},
			want:    nil,
			wantErr: zeroerr.ErrCategoryNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockProduct := NewProductService(t)
			tt.mock(mockProduct)

			svcCtx := &svc.ServiceContext{
				ProductRpc: mockProduct,
			}

			logic := NewCreateProductLogic(context.Background(), svcCtx)
			got, err := logic.CreateProduct(tt.req)

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
