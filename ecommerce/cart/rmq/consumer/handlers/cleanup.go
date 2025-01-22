package handlers

import (
	"context"
	"fmt"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/middleware"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
)

type CleanupHandler struct {
    BaseHandler
    cartItemModel      model.CartItemsModel
    cartStatisticsModel model.CartStatisticsModel
}

func NewCleanupHandler(
    logger middleware.Logger,
    cartItemModel model.CartItemsModel,
    cartStatisticsModel model.CartStatisticsModel,
) *CleanupHandler {
    return &CleanupHandler{
        BaseHandler:         NewBaseHandler(logger),
        cartItemModel:      cartItemModel,
        cartStatisticsModel: cartStatisticsModel,
    }
}

func (h *CleanupHandler) Handle(event *types.CartEvent) error {
    h.LogEvent(event, "handling cleanup event")

    data, ok := event.Data.(*types.CartClearedData)
    if !ok {
        return fmt.Errorf("invalid event data type")
    }

    // Delete all cart items for user
    if err := h.cartItemModel.DeleteByUserId(context.Background(), uint64(data.UserID)); err != nil {
        return err
    }

    // Reset cart statistics
    if err := h.cartStatisticsModel.Delete(context.Background(), uint64(data.UserID)); err != nil {
        return err
    }

    h.LogEvent(event, "cleanup handled successfully")
    return nil
}