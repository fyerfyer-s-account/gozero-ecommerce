package product

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/application/internal/types"
	product "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

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

func (l *GetProductReviewsLogic) GetProductReviews(req *types.ReviewListReq) ([]types.Review, error) {
	// Call product RPC
	resp, err := l.svcCtx.ProductRpc.ListReviews(l.ctx, &product.ListReviewsRequest{
		ProductId: req.ProductId,
		Page:      req.Page,
	})
	if err != nil {
		return nil, err
	}

	// Convert reviews
	reviews := make([]types.Review, 0, len(resp.Reviews))
	for _, r := range resp.Reviews {
		reviews = append(reviews, types.Review{
			Id:        r.Id,
			ProductId: r.ProductId,
			OrderId:   r.OrderId,
			UserId:    r.UserId,
			Rating:    r.Rating,
			Content:   r.Content,
			Images:    r.Images,
			CreatedAt: r.CreatedAt,
		})
	}

	return reviews, nil
}
