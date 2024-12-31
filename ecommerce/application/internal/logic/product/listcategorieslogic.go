package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewListCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategoriesLogic {
	return &ListCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ListCategoriesLogic) ListCategories() ([]types.Category, error) {
	resp, err := l.svcCtx.ProductRpc.ListCategories(l.ctx, &product.ListCategoriesRequest{
		ParentId: 0, // Get root categories by default
	})
	if err != nil {
		return nil, err
	}

	categories := make([]types.Category, 0, len(resp.Categories))
	for _, c := range resp.Categories {
		categories = append(categories, types.Category{
			Id:       c.Id,
			Name:     c.Name,
			ParentId: c.ParentId,
			Level:    int32(c.Level),
			Sort:     int32(c.Sort),
			Icon:     c.Icon,
		})
	}

	return categories, nil
}
