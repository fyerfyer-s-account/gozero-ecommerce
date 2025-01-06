package logic

import (
	"context"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
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
	// Validate input
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check if review exists
	_, err := l.svcCtx.ProductReviewsModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrReviewNotFound
	}

	// Update review Status
	err = l.svcCtx.ProductReviewsModel.UpdateStatus(l.ctx, uint64(in.Id), in.Status)
	if err != nil {
		return nil, err
	}

	return &product.UpdateReviewResponse{Success: true}, nil
}
