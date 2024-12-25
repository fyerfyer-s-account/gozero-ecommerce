package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type CreateReviewLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewCreateReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CreateReviewLogic {
	return &CreateReviewLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 评价管理
func (l *CreateReviewLogic) CreateReview(in *product.CreateReviewRequest) (*product.CreateReviewResponse, error) {
	// todo: add your logic here and delete this line

	return &product.CreateReviewResponse{}, nil
}
