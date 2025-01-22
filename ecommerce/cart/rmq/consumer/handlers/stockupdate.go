package handlers

import (
    "context"
    "fmt"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/middleware"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
)

type StockUpdateHandler struct {
    BaseHandler
    cartItemModel      model.CartItemsModel
    cartStatisticsModel model.CartStatisticsModel
}

func NewStockUpdateHandler(
    logger middleware.Logger,
    cartItemModel model.CartItemsModel,
    cartStatisticsModel model.CartStatisticsModel,
) *StockUpdateHandler {
    return &StockUpdateHandler{
        BaseHandler:         NewBaseHandler(logger),
        cartItemModel:      cartItemModel,
        cartStatisticsModel: cartStatisticsModel,
    }
}

func (h *StockUpdateHandler) Handle(event *types.CartEvent) error {
    h.LogEvent(event, "handling stock update event")

    data, ok := event.Data.(*types.CartItemData)
    if !ok {
        return fmt.Errorf("invalid event data type")
    }

    // Update cart item quantity if necessary
    item, err := h.cartItemModel.FindOneByUserIdSkuId(context.Background(), uint64(data.UserID), uint64(data.ProductID))
    if err != nil {
        return err
    }

    if item.Quantity > int64(data.Quantity) {
        if err := h.cartItemModel.UpdateQuantity(
            context.Background(),
            uint64(data.UserID),
            uint64(data.ProductID),
            uint64(data.ProductID),
            int64(data.Quantity),
        ); err != nil {
            return err
        }

        // Recalculate cart statistics
        if err := h.cartStatisticsModel.RecalculateStats(
            context.Background(),
            uint64(data.UserID),
        ); err != nil {
            return err
        }
    }

    h.LogEvent(event, "stock update handled successfully")
    return nil
}