package logic

import (
	"context"
	"encoding/json"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventory"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"

	"github.com/zeromicro/go-zero/core/logx"
)

type GetSkuLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetSkuLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetSkuLogic {
	return &GetSkuLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetSkuLogic) GetSku(in *product.GetSkuRequest) (*product.GetSkuResponse, error) {
	// Validate input
	if in.Id <= 0 {
		return nil, zeroerr.ErrInvalidParam
	}

	// Get SKU
	sku, err := l.svcCtx.SkusModel.FindOne(l.ctx, uint64(in.Id))
	if err != nil {
		return nil, zeroerr.ErrSkuNotFound
	}

	inventory, err := l.svcCtx.InventoryRpc.GetStock(l.ctx, &inventory.GetStockRequest{
		SkuId: int64(sku.Id),
		WarehouseId: l.svcCtx.Config.DefaultWarehouseId,
	})

	if err != nil {
		return nil, zeroerr.ErrProductNotFound
	}

	// Convert to proto message
	pbSku := &product.Sku{
		Id:        int64(sku.Id),
		ProductId: int64(sku.ProductId),
		SkuCode:   sku.SkuCode,
		Price:     sku.Price,
		Stock:     int64(inventory.Stock.Available),
		CreatedAt: sku.CreatedAt.Unix(),
		UpdatedAt: sku.UpdatedAt.Unix(),
	}

	// Parse attributes JSON
	if sku.Attributes != "" {
		var attrs []*product.SkuAttribute
		if err := json.Unmarshal([]byte(sku.Attributes), &attrs); err == nil {
			pbSku.Attributes = attrs
		}
	}

	return &product.GetSkuResponse{
		Sku: pbSku,
	}, nil
}
