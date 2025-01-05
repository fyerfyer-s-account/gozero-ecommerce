package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/productservice"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetCategoriesLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetCategoriesLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetCategoriesLogic {
	return &GetCategoriesLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetCategoriesLogic) GetCategories() (resp *types.GetCategoriesResp, err error) {
	// todo: add your logic here and delete this line
	res, err := l.svcCtx.ProductRpc.GetCategories(l.ctx, &productservice.Empty{})
	if err != nil {
		return nil, err 
	}

	return &types.GetCategoriesResp{
		CategoryNames: res.CategoryNames,
	}, nil
}
