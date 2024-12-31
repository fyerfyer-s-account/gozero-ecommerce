package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeReviewStatusLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewChangeReviewStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeReviewStatusLogic {
	return &ChangeReviewStatusLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ChangeReviewStatusLogic) ChangeReviewStatus(req *types.ChangeReviewStatusReq) error {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.ProductRpc.ChangeReviewStatus(l.ctx, &product.ChangeReviewStatusRequest{
		Id:     req.Id,
		Status: int64(req.Status),
	})
	if err != nil {
		return err
	}

	return nil
}
