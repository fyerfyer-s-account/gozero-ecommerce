package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListCategoriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListCategoriesLogic {
	return &ListCategoriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListCategoriesLogic) ListCategories(in *product.ListCategoriesRequest) (*product.ListCategoriesResponse, error) {
	// Get categories by parent ID
	categories, err := l.svcCtx.CategoriesModel.FindByParentId(l.ctx, uint64(in.ParentId))
	if err != nil {
		logx.Errorf("Failed to get categories: %v", err)
		return nil, zeroerr.ErrCategoryNotFound
	}

	// Convert to proto message
	var pbCategories []*product.Category
	for _, category := range categories {
		pbCategory := &product.Category{
			Id:       int64(category.Id),
			Name:     category.Name,
			ParentId: int64(category.ParentId.Int64),
			Level:    category.Level,
			Sort:     category.Sort,
			Icon:     category.Icon.String,
		}
		pbCategories = append(pbCategories, pbCategory)
	}

	return &product.ListCategoriesResponse{
		Categories: pbCategories,
	}, nil
}
