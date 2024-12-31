package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewDeleteReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteReviewLogic {
	return &DeleteReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DeleteReviewLogic) DeleteReview(req *types.DeleteReviewReq) error {
	// todo: add your logic here and delete this line
	_, err := l.svcCtx.ProductRpc.DeleteReview(l.ctx, &product.DeleteReviewRequest{
		Id: req.Id,
	})
	return err
}
