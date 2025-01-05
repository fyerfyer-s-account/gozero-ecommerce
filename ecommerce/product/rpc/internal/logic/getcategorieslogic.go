package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoriesLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoriesLogic {
	return &GetCategoriesLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetCategoriesLogic) GetCategories(in *product.Empty) (*product.GetCategoriesResponse, error) {
	// todo: add your logic here and delete this line

	cats, err := l.svcCtx.CategoriesModel.GetAll(l.ctx);
	if err != nil {
		return nil, err 
	}

	var resp []string
	for _, cat := range cats {
		resp = append(resp, cat.Name)
	}

	return &product.GetCategoriesResponse{CategoryNames: resp}, nil
}
