package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateCategoryLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCreateCategoryLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateCategoryLogic {
	return &CreateCategoryLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CreateCategoryLogic) CreateCategory(req *types.CreateCategoryReq) (*types.CreateCategoryResp, error) {
	resp, err := l.svcCtx.ProductRpc.CreateCategory(l.ctx, &product.CreateCategoryRequest{
		Name:     req.Name,
		ParentId: req.ParentId,
		Sort:     int64(req.Sort),
		Icon:     req.Icon,
	})
	if err != nil {
		return nil, err
	}

	return &types.CreateCategoryResp{
		Id: resp.Id,
	}, nil
}
