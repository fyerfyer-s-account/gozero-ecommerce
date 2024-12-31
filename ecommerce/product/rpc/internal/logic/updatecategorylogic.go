package logic

import (
	"context"
	"database/sql"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateCategoryLogic {
	return &UpdateCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateCategoryLogic) UpdateCategory(in *product.UpdateCategoryRequest) (*product.UpdateCategoryResponse, error) {
    // Validate input
    if in.Id <= 0 || in.Name == "" {
        return nil, zeroerr.ErrInvalidParam
    }

    // Check category exists
    category, err := l.svcCtx.CategoriesModel.FindOne(l.ctx, uint64(in.Id))
    if err != nil {
        return nil, zeroerr.ErrCategoryNotFound
    }

    // Check name uniqueness if changed
    if category.Name != in.Name {
        existing, err := l.svcCtx.CategoriesModel.FindOneByName(l.ctx, in.Name)
        if err != nil && err != model.ErrNotFound {
            return nil, err
        }
        if existing != nil && existing.Id != category.Id {
            return nil, zeroerr.ErrCategoryDuplicate
        }
    }

    // Update category
    category.Name = in.Name
    category.Sort = in.Sort
    if in.Icon != "" {
        category.Icon = sql.NullString{String: in.Icon, Valid: true}
    }

    err = l.svcCtx.CategoriesModel.Update(l.ctx, category)
    if err != nil {
        logx.Errorf("Failed to update category: %v", err)
        return nil, zeroerr.ErrCategoryUpdateFailed
    }

    return &product.UpdateCategoryResponse{
        Success: true,
    }, nil
}
