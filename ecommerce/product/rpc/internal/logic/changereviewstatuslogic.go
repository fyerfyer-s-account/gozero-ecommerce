package logic

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ChangeReviewStatusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewChangeReviewStatusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ChangeReviewStatusLogic {
	return &ChangeReviewStatusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ChangeReviewStatusLogic) ChangeReviewStatus(in *product.ChangeReviewStatusRequest) (*product.ChangeReviewStatusResponse, error) {
	// Validate input
	if in.Id <= 0 || in.Status < 0 || in.Status > 2 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Check review exists
	_, err := l.svcCtx.ProductReviewsModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrReviewNotFound
	}

	// Update status
	err = l.svcCtx.ProductReviewsModel.UpdateStatus(l.ctx, uint64(in.Id), in.Status)
	if err != nil {
		logx.Errorf("Failed to update review status: %v", err)
		return nil, zeroerr.ErrReviewUpdateFailed
	}

	return &product.ChangeReviewStatusResponse{
		Success: true,
	}, nil
}
