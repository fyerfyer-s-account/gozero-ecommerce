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

type CreateCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
	return &CreateCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 分类管理
func (l *CreateCategoryLogic) CreateCategory(in *product.CreateCategoryRequest) (*product.CreateCategoryResponse, error) {
    // Validate input
    if in.Name == "" {
        return nil, zeroerr.ErrInvalidParam
    }

    // Check name uniqueness
    existingCat, err := l.svcCtx.CategoriesModel.FindOneByName(l.ctx, in.Name)
    if err != nil && err != model.ErrNotFound {
        return nil, err
    }
    if existingCat != nil {
        return nil, zeroerr.ErrCategoryDuplicate
    }

    // Get category level
    level, err := l.svcCtx.CategoriesModel.GetLevel(l.ctx, uint64(in.ParentId))
    if err != nil && in.ParentId != 0 {
        return nil, zeroerr.ErrCategoryNotFound
    }

    // Check max level
    if level > int64(l.svcCtx.Config.MaxCategoryLevel) {
        return nil, zeroerr.ErrInvalidCategoryLevel
    }

    // Create category
    category := &model.Categories{
        Name:     in.Name,
        ParentId: sql.NullInt64{Int64: in.ParentId, Valid: in.ParentId != 0},
        Level:    level,
        Sort:     in.Sort,
        Icon:     sql.NullString{String: in.Icon, Valid: in.Icon != ""},
    }

    result, err := l.svcCtx.CategoriesModel.Insert(l.ctx, category)
    if err != nil {
        logx.Errorf("Failed to create category: %v", err)
        return nil, zeroerr.ErrCategoryCreateFailed
    }

    id, err := result.LastInsertId()
    if err != nil {
        return nil, err
    }

    return &product.CreateCategoryResponse{
        Id: id,
    }, nil
}
