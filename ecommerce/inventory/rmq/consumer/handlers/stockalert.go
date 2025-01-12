package handlers

import (
	"context"
	"fmt"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
)

type StockAlertHandler struct {
    inventoryRpc inventoryclient.Inventory
    messageRpc   messageservice.MessageService
}

func NewStockAlertHandler(inventoryRpc inventoryclient.Inventory, messageRpc messageservice.MessageService) *StockAlertHandler {
    return &StockAlertHandler{
        inventoryRpc: inventoryRpc,
        messageRpc:   messageRpc,
    }
}

func (h *StockAlertHandler) Handle(event *types.InventoryEvent) error {
    data, ok := event.Data.(*types.StockAlertData)
    if !ok {
        return fmt.Errorf("invalid event data type")
    }

    if data.Available <= data.Threshold {
        _, err := h.messageRpc.SendTemplateMessage(context.Background(), &messageservice.SendTemplateMessageRequest{
            TemplateCode: "stock_alert",
            UserId:      event.Metadata.UserID,
            Params: map[string]string{
                "skuId":       fmt.Sprint(data.SkuID),
                "warehouse":   fmt.Sprint(data.WarehouseID),
                "available":   fmt.Sprint(data.Available),
                "threshold":   fmt.Sprint(data.Threshold),
            },
            Channels: []int32{1, 2}, 
        })
        return err
    }

    return nil
}