package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
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
	// Validate input
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check review exists
	_, err := l.svcCtx.ProductReviewsModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrReviewNotFound
	}

	// Delete review
	err = l.svcCtx.ProductReviewsModel.Delete(l.ctx, uint64(in.Id))
	if err != nil {
		logx.Errorf("Failed to delete review %d: %v", in.Id, err)
		return nil, zeroerr.ErrReviewDeleteFailed
	}

	return &product.DeleteReviewResponse{
		Success: true,
	}, nil
}
