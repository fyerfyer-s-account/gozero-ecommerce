package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoryLogic {
	return &GetCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCategoryLogic) GetCategory(in *product.GetCategoryRequest) (*product.GetCategoryResponse, error) {
	// Validate input
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Get category
	category, err := l.svcCtx.CategoriesModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrCategoryNotFound
	}

	// Get child categories
	children, err := l.svcCtx.CategoriesModel.FindByParentId(l.ctx, uint64(in.Id))
	if err != nil {
		logx.Errorf("Failed to get child categories for category %d: %v", in.Id, err)
		return nil, err
	}

	// Convert to proto message
	pbCategory := &product.Category{
		Id:       int64(category.Id),
		Name:     category.Name,
		ParentId: int64(category.ParentId.Int64),
		Level:    category.Level,
		Sort:     category.Sort,
		Icon:     category.Icon.String,
	}

	// Convert children to proto messages
	pbChildren := make([]*product.Category, 0, len(children))
	for _, child := range children {
		pbChild := &product.Category{
			Id:       int64(child.Id),
			Name:     child.Name,
			ParentId: int64(child.ParentId.Int64),
			Level:    child.Level,
			Sort:     child.Sort,
			Icon:     child.Icon.String,
		}
		pbChildren = append(pbChildren, pbChild)
	}

	return &product.GetCategoryResponse{
		Category: pbCategory,
		Children: pbChildren,
	}, nil
}
