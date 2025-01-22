package logic

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
	"github.com/google/uuid"

	"github.com/zeromicro/go-zero/core/logx"
)

type AddItemLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewAddItemLogic(ctx context.Context, svcCtx *svc.ServiceContext) *AddItemLogic {
	return &AddItemLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

// 购物车操作
func (l *AddItemLogic) AddItem(in *cart.AddItemRequest) (*cart.AddItemResponse, error) {
    if in.UserId <= 0 || in.ProductId <= 0 || in.SkuId <= 0 || in.Quantity <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Get product and SKU info
    sku, err := l.svcCtx.ProductRpc.GetSku(l.ctx, &product.GetSkuRequest{
        Id: in.SkuId,
    })
    if err != nil {
        return nil, err
    }

    if sku.Sku.Stock < int64(in.Quantity) {
        return nil, zeroerr.ErrItemOutOfStock
    }

    prod, err := l.svcCtx.ProductRpc.GetProduct(l.ctx, &product.GetProductRequest{
        Id: in.ProductId,
    })
    if err != nil {
        return nil, err
    }

    // Check existing cart item
    cartItem, err := l.svcCtx.CartItemsModel.FindOneByUserIdSkuId(l.ctx, uint64(in.UserId), uint64(in.SkuId))
    if err != nil && err != model.ErrNotFound {
        return nil, err
    }

    if cartItem != nil {
        // Update existing item
        err = l.svcCtx.CartItemsModel.UpdateQuantity(l.ctx, uint64(in.UserId), uint64(in.ProductId), uint64(in.SkuId), cartItem.Quantity+in.Quantity)
        if err != nil {
            return nil, err
        }
    } else {
        // Create new item
        cartItem = &model.CartItems{
            UserId:      uint64(in.UserId),
            ProductId:   uint64(in.ProductId),
            SkuId:      uint64(in.SkuId),
            ProductName: prod.Product.Name,
            SkuName:    sku.Sku.SkuCode,
            Image: sql.NullString{
                String: prod.Product.Images[0],
                Valid:  prod.Product.Images[0] != "",
            },
            Price:    sku.Sku.Price,
            Quantity: int64(in.Quantity),
            Selected: 1,
        }
        _, err = l.svcCtx.CartItemsModel.Insert(l.ctx, cartItem)
        if err != nil {
            return nil, err
        }
    }

    // Update cart statistics
    stats, err := l.svcCtx.CartStatsModel.FindOne(l.ctx, uint64(in.UserId))
    if err == model.ErrNotFound {
        stats = &model.CartStatistics{
            UserId:           uint64(in.UserId),
            TotalQuantity:    int64(in.Quantity),
            SelectedQuantity: int64(in.Quantity),
            TotalAmount:      sku.Sku.Price * float64(in.Quantity),
            SelectedAmount:   sku.Sku.Price * float64(in.Quantity),
        }
        err = l.svcCtx.CartStatsModel.Upsert(l.ctx, stats)
    } else if err == nil {
        err = l.svcCtx.CartStatsModel.RecalculateStats(l.ctx, uint64(in.UserId))
    }
    if err != nil {
        return nil, err
    }
    
    event := &types.CartEvent{
        ID:        uuid.New().String(),
        Type:      types.EventTypeItemAdded,
        Timestamp: time.Now(),
        Data: &types.CartItemAddedData{
            CartItemData: types.CartItemData{
                UserID:    in.UserId,
                ProductID: in.ProductId,
                Quantity:  int32(in.Quantity),
            },
        },
        Metadata: types.EventMetadata{
            TraceID: l.ctx.Value("trace_id").(string),
            UserID:  fmt.Sprint(in.UserId),
        },
    }

    if err := l.svcCtx.Producer.PublishEvent(event); err != nil {
        logx.Errorf("Failed to publish cart.item.added event: %v", err)
        // Continue execution as this is not critical
    }

    return &cart.AddItemResponse{
        Success: true,
    }, nil
}
