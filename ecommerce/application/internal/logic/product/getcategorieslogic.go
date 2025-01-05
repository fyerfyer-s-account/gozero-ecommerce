package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

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
	res, err := l.svcCtx.ProductRpc.GetCategories(l.ctx, &product.Empty{})
	if err != nil {
		return nil, err 
	}

	var cats []types.Category
	for _, cat := range res.GetCategories() {
		cats = append(cats, types.Category {
			Id: cat.Id,
			Name: cat.Name,
			ParentId: cat.ParentId,
			Level: int32(cat.Level),
			Sort: int32(cat.Sort),
			Icon: cat.Icon,
		})
	}

	return &types.GetCategoriesResp{
		Categories: cats,
	}, nil
}
