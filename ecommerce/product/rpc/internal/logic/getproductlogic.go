package logic

import (
	"context"
	"encoding/json"
	"strings"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetProductLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductLogic {
	return &GetProductLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetProductLogic) GetProduct(in *product.GetProductRequest) (*product.GetProductResponse, error) {
	// Validate input
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Get product
	prod, err := l.svcCtx.ProductsModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrProductNotFound
	}

	// Get SKUs
	skus, err := l.svcCtx.SkusModel.FindManyByProductId(l.ctx, uint64(in.Id))
	if err != nil {
		logx.Errorf("Failed to get SKUs for product %d: %v", in.Id, err)
		return nil, err
	}

	// Convert product to proto message
	pbProduct := &product.Product{
		Id:          int64(prod.Id),
		Name:        prod.Name,
		Description: prod.Description.String,
		CategoryId:  int64(prod.CategoryId),
		Brand:       prod.Brand.String,
		Price:       prod.Price,
		Status:      prod.Status,
		CreatedAt:   prod.CreatedAt.Unix(),
		UpdatedAt:   prod.UpdatedAt.Unix(),
	}

	if prod.Images.Valid {
		pbProduct.Images = strings.Split(prod.Images.String, ",")
	}

	// Convert SKUs to proto messages
	pbSkus := make([]*product.Sku, 0, len(skus))
	for _, sku := range skus {
		pbSku := &product.Sku{
			Id:        int64(sku.Id),
			ProductId: int64(sku.ProductId),
			SkuCode:   sku.SkuCode,
			Price:     sku.Price,
			Stock:     sku.Stock,
			CreatedAt: sku.CreatedAt.Unix(),
			UpdatedAt: sku.UpdatedAt.Unix(),
		}

		// Parse attributes JSON
		if sku.Attributes != "" {
			var attrs []*product.SkuAttribute
			if err := json.Unmarshal([]byte(sku.Attributes), &attrs); err == nil {
				pbSku.Attributes = attrs
			}
		}

		pbSkus = append(pbSkus, pbSku)
	}

	return &product.GetProductResponse{
		Product: pbProduct,
		Skus:    pbSkus,
	}, nil
}
