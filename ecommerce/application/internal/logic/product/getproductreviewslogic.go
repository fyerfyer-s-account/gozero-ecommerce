package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetProductReviewsLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewGetProductReviewsLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetProductReviewsLogic {
	return &GetProductReviewsLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *GetProductReviewsLogic) GetProductReviews(req *types.ReviewListReq) (resp []types.Review, err error) {
	// todo: add your logic here and delete this line

	return
}
