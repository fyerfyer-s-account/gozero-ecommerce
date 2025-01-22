package handlers

import (
    "context"
    "fmt"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/middleware"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
)

type PriceUpdateHandler struct {
    BaseHandler
    cartItemModel      model.CartItemsModel
    cartStatisticsModel model.CartStatisticsModel
}

func NewPriceUpdateHandler(
    logger middleware.Logger,
    cartItemModel model.CartItemsModel,
    cartStatisticsModel model.CartStatisticsModel,
) *PriceUpdateHandler {
    return &PriceUpdateHandler{
        BaseHandler:         NewBaseHandler(logger),
        cartItemModel:      cartItemModel,
        cartStatisticsModel: cartStatisticsModel,
    }
}

func (h *PriceUpdateHandler) Handle(event *types.CartEvent) error {
    h.LogEvent(event, "handling price update event")

    data, ok := event.Data.(*types.CartItemData)
    if !ok {
        return fmt.Errorf("invalid event data type")
    }

    // Recalculate cart statistics after price update
    if err := h.cartStatisticsModel.RecalculateStats(
        context.Background(),
        uint64(data.UserID),
    ); err != nil {
        return err
    }

    h.LogEvent(event, "price update handled successfully")
    return nil
}