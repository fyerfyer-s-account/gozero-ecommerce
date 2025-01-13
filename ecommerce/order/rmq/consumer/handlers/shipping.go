package handlers

import (
	"context"
	"fmt"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
)

type ShippingHandler struct {
	messageRpc         messageservice.MessageService
	orderModel         model.OrdersModel
	orderShippingModel model.OrderShippingModel
}

func NewShippingHandler(messageRpc messageservice.MessageService, orderModel model.OrdersModel,
	orderShippingModel model.OrderShippingModel) *ShippingHandler {
	return &ShippingHandler{
		messageRpc:         messageRpc,
		orderModel:         orderModel,
		orderShippingModel: orderShippingModel,
	}
}

func (h *ShippingHandler) Handle(event *types.OrderEvent) error {
	switch event.Type {
	case types.EventTypeOrderShipped:
		return h.handleOrderShipped(event)
	default:
		return nil
	}
}

func (h *ShippingHandler) handleOrderShipped(event *types.OrderEvent) error {
	data, ok := event.Data.(*types.OrderShippedData)
	if !ok {
		return fmt.Errorf("invalid event data type")
	}

	// 1. Find order shipping record
	shipping, err := h.orderShippingModel.FindByOrderId(context.Background(), uint64(data.OrderId))
	if err != nil {
		return err
	}

	// 2. Update shipping info
	err = h.orderShippingModel.UpdateShippingInfo(context.Background(),
		shipping.OrderId, data.ShippingNo, data.Company)
	if err != nil {
		return err
	}

	// 3. Send notification
	_, err = h.messageRpc.SendTemplateMessage(context.Background(), &messageservice.SendTemplateMessageRequest{
		TemplateCode: "order_shipped",
		UserId:       event.Metadata.UserID,
		Params: map[string]string{
			"orderNo":    data.OrderNo,
			"shippingNo": data.ShippingNo,
			"company":    data.Company,
		},
	})

	return err
}
