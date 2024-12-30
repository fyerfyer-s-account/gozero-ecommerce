package logic

import (
	"context"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
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
	if in.ProductId <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	if in.PageSize <= 0 {
		in.PageSize = 10
	}
	if in.Page <= 0 {
		in.Page = 1
	}

	// Get total count
	total, err := l.svcCtx.ProductReviewsModel.Count(l.ctx, uint64(in.ProductId))
	if err != nil {
		logx.Errorf("Failed to get reviews count: %v", err)
		return nil, err
	}

	// Get reviews with filters
	reviews, err := l.svcCtx.ProductReviewsModel.FindManyByProductId(
		l.ctx,
		uint64(in.ProductId),
		int(in.Page),
		int(in.PageSize),
	)
	if err != nil {
		logx.Errorf("Failed to get reviews: %v", err)
		return nil, err
	}

	// Convert to proto messages
	pbReviews := make([]*product.Review, 0, len(reviews))
	for _, r := range reviews {
		pbReview := &product.Review{
			Id:        int64(r.Id),
			ProductId: int64(r.ProductId),
			UserId:    int64(r.UserId),
			OrderId:   int64(r.OrderId),
			Rating:    int32(r.Rating),
			Content:   r.Content.String,
			Status:    r.Status,
			CreatedAt: r.CreatedAt.Unix(),
			UpdatedAt: r.UpdatedAt.Unix(),
		}

		if r.Images.Valid {
			var images []string
			if err := json.Unmarshal([]byte(r.Images.String), &images); err == nil {
				pbReview.Images = images
			}
		}
		pbReviews = append(pbReviews, pbReview)
	}

	return &product.ListReviewsResponse{
		Total:   total,
		Reviews: pbReviews,
	}, nil
}
