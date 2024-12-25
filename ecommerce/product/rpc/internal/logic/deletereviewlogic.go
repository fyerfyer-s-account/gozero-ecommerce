package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type DeleteReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewDeleteReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DeleteReviewLogic {
	return &DeleteReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *DeleteReviewLogic) DeleteReview(in *product.DeleteReviewRequest) (*product.DeleteReviewResponse, error) {
	// todo: add your logic here and delete this line

	return &product.DeleteReviewResponse{}, nil
}
