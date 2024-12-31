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
	if in.Id <= 0 || in.Rating < 1 || in.Rating > 5 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check review exists
	_, err := l.svcCtx.ProductReviewsModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrReviewNotFound
	}

	// Convert images to JSON
	var imagesJSON string
	if len(in.Images) > 0 {
		imagesBytes, err := json.Marshal(in.Images)
		if err != nil {
			return nil, zeroerr.ErrInvalidParam
		}
		imagesJSON = string(imagesBytes)
	}

	// Update only allowed fields
	err = l.svcCtx.ProductReviewsModel.UpdateContent(
		l.ctx,
		uint64(in.Id),
		int64(in.Rating),
		sql.NullString{String: in.Content, Valid: in.Content != ""},
		sql.NullString{String: imagesJSON, Valid: len(in.Images) > 0},
	)
	if err != nil {
		logx.Errorf("Failed to update review: %v", err)
		return nil, zeroerr.ErrReviewUpdateFailed
	}

	return &product.UpdateReviewResponse{
		Success: true,
	}, nil
}
