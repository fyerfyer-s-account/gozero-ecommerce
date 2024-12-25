package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListReviewsLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListReviewsLogic {
	return &ListReviewsLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListReviewsLogic) ListReviews(in *product.ListReviewsRequest) (*product.ListReviewsResponse, error) {
	// todo: add your logic here and delete this line

	return &product.ListReviewsResponse{}, nil
}
