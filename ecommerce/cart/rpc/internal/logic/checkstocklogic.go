package logic

import (
    "context"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/cart"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zeroerr"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/product/rpc/product"
    "github.com/zeromicro/go-zero/core/logx"
)

type CheckStockLogic struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logx.Logger
}

func NewCheckStockLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CheckStockLogic {
    return &CheckStockLogic{
        ctx:    ctx,
        svcCtx: svcCtx,
        Logger: logx.WithContext(ctx),
    }
}

func (l *CheckStockLogic) CheckStock(in *cart.CheckStockRequest) (*cart.CheckStockResponse, error) {
    if in.UserId <= 0 {
        return nil, zeroerr.ErrInvalidParam
    }

    // Get selected items from cart
    items, err := l.svcCtx.CartItemsModel.FindSelectedByUserId(l.ctx, uint64(in.UserId))
    if err != nil {
        return nil, err
    }

    if len(items) == 0 {
        return &cart.CheckStockResponse{
            AllInStock: true,
        }, nil
    }

    outOfStockItems := make([]*cart.CartItem, 0)

    // Check stock for each item
    for _, item := range items {
        // Get SKU info from product service
        skuResp, err := l.svcCtx.ProductRpc.GetSku(l.ctx, &product.GetSkuRequest{
            Id: int64(item.SkuId),
        })
        if err != nil {
            return nil, err
        }

        if skuResp.Sku.Stock < item.Quantity {
            outOfStockItems = append(outOfStockItems, &cart.CartItem{
                Id:          int64(item.Id),
                UserId:      int64(item.UserId),
                ProductId:   int64(item.ProductId),
                SkuId:       int64(item.SkuId),
                ProductName: item.ProductName,
                SkuName:     item.SkuName,
                Image:       item.Image.String,
                Price:       item.Price,
                Quantity:    item.Quantity,
                Selected:    item.Selected == 1,
                Stock:       int32(skuResp.Sku.Stock),
                CreatedAt:   item.CreatedAt.Unix(),
                UpdatedAt:   item.UpdatedAt.Unix(),
            })
        }
    }

    return &cart.CheckStockResponse{
        AllInStock:       len(outOfStockItems) == 0,
        OutOfStockItems: outOfStockItems,
    }, nil
}