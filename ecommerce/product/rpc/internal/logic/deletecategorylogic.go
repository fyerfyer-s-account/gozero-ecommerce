package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteCategoryLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteCategoryLogic {
	return &DeleteCategoryLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteCategoryLogic) DeleteCategory(in *product.DeleteCategoryRequest) (*product.DeleteCategoryResponse, error) {
	// Validate input
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check category exists
	_, err := l.svcCtx.CategoriesModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrCategoryNotFound
	}

	// Check if category has children
	hasChildren, err := l.svcCtx.CategoriesModel.HasChildren(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, err
	}
	if hasChildren {
		return nil, zeroerr.ErrCategoryHasChildren
	}

	// Check if category has products
	hasProducts, err := l.svcCtx.CategoriesModel.HasProducts(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, err
	}
	if hasProducts {
		return nil, zeroerr.ErrCategoryHasProducts
	}

	// Delete category
	err = l.svcCtx.CategoriesModel.Delete(l.ctx, uint64(in.Id))
	if err != nil {
		logx.Errorf("Failed to delete category %d: %v", in.Id, err)
		return nil, zeroerr.ErrCategoryDeleteFailed
	}

	return &product.DeleteCategoryResponse{
		Success: true,
	}, nil
}
