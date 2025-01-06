package logic

import (
	"context"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type ListSkusLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewListSkusLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ListSkusLogic {
	return &ListSkusLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *ListSkusLogic) ListSkus(in *product.ListSkusRequest) (*product.ListSkusResponse, error) {
	if in.ProductId <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	if in.Page <= 0 {
		in.Page = 1
	}

	// Get total count
	total, err := l.svcCtx.SkusModel.Count(l.ctx, uint64(in.ProductId))
	if err != nil {
		logx.Errorf("Failed to get SKUs count: %v", err)
		return nil, err
	}

	// Get SKUs with pagination
	skus, err := l.svcCtx.SkusModel.FindManyPageByProductId(
		l.ctx,
		uint64(in.ProductId),
		int(in.Page),
		l.svcCtx.Config.PageSize,
	)
	if err != nil {
		logx.Errorf("Failed to get SKUs: %v", err)
		return nil, err
	}

	// Convert to proto messages
	pbSkus := make([]*product.Sku, 0, len(skus))
	for _, s := range skus {
		pbSku := &product.Sku{
			Id:        int64(s.Id),
			ProductId: int64(s.ProductId),
			SkuCode:   s.SkuCode,
			Price:     s.Price,
			Stock:     s.Stock,
			Sales:     s.Sales,
		}

		// Parse attributes JSON
		if s.Attributes.Valid {
			var attrs []*product.SkuAttribute
			if err := json.Unmarshal([]byte(s.Attributes.String), &attrs); err == nil {
				pbSku.Attributes = attrs
			}
		}

		pbSkus = append(pbSkus, pbSku)
	}

	return &product.ListSkusResponse{
		Total: total,
		Skus:  pbSkus,
	}, nil
}
