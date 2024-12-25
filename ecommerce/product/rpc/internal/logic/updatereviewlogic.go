package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewUpdateReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReviewLogic {
	return &UpdateReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *UpdateReviewLogic) UpdateReview(in *product.UpdateReviewRequest) (*product.UpdateReviewResponse, error) {
	// todo: add your logic here and delete this line

	return &product.UpdateReviewResponse{}, nil
}
