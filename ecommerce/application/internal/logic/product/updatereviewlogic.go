package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type UpdateReviewLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUpdateReviewLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UpdateReviewLogic {
	return &UpdateReviewLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UpdateReviewLogic) UpdateReview(req *types.UpdateReviewReq) error {
	// todo: add your logic here and delete this line

	return nil
}
