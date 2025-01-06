package logic

import (
	"context"
	"database/sql"
	"encoding/json"
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

	// Prepare updates
	updates := make(map[string]interface{})
	if in.Rating > 0 {
		updates["rating"] = in.Rating
	}
	if in.Content != "" {
		updates["content"] = sql.NullString{String: in.Content, Valid: true}
	}
	if len(in.Images) > 0 {
		imagesJSON, err := json.Marshal(in.Images)
		if err != nil {
			return nil, err
		}
		updates["images"] = sql.NullString{String: string(imagesJSON), Valid: true}
	}
	if in.Status > 0 {
		updates["status"] = in.Status
	}

	// Update review
	err = l.svcCtx.ProductReviewsModel.UpdateReviews(l.ctx, uint64(in.Id), updates)
	if err != nil {
		return nil, err
	}

	return &product.UpdateReviewResponse{Success: true}, nil
}
