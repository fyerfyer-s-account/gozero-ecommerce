package logic

import (
	"context"
	"database/sql"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/model"
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
	// Validate input
	if in.ProductId <= 0 || in.UserId <= 0 || in.OrderId <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	if in.Rating < 1 || in.Rating > 5 {
		return nil, zeroerr.ErrInvalidParam
	}

	contentLen := len(in.Content)
	if contentLen < l.svcCtx.Config.MinReviewLength || contentLen > l.svcCtx.Config.MaxReviewLength {
		return nil, zeroerr.ErrInvalidParam
	}

	if len(in.Images) > l.svcCtx.Config.MaxReviewImages {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check product exists
	_, err := l.svcCtx.ProductsModel.FindOne(l.ctx, uint64(in.ProductId))
	if err != nil {
		return nil, zeroerr.ErrProductNotFound
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

	// Create review
	review := &model.ProductReviews{
		ProductId: uint64(in.ProductId),
		UserId:    uint64(in.UserId),
		OrderId:   uint64(in.OrderId),
		Rating:    int64(in.Rating),
		Content:   sql.NullString{String: in.Content, Valid: true},
		Images:    sql.NullString{String: imagesJSON, Valid: len(in.Images) > 0},
		Status:    0, // Default to pending review
	}

	result, err := l.svcCtx.ProductReviewsModel.Insert(l.ctx, review)
	if err != nil {
		logx.Errorf("Failed to create review: %v", err)
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}

	return &product.CreateReviewResponse{
		Id: id,
	}, nil
}
