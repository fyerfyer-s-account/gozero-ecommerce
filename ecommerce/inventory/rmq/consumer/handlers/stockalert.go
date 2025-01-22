package handlers

import (
    "context"
    "fmt"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/middleware"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
)

type StockAlertHandler struct {
    BaseHandler
    inventoryRpc inventoryclient.Inventory
    messageRpc   messageservice.MessageService
}

func NewStockAlertHandler(
    logger middleware.Logger,
    inventoryRpc inventoryclient.Inventory,
    messageRpc messageservice.MessageService,
) *StockAlertHandler {
    return &StockAlertHandler{
        BaseHandler:   NewBaseHandler(logger),
        inventoryRpc: inventoryRpc,
        messageRpc:   messageRpc,
    }
}

func (h *StockAlertHandler) Handle(event *types.InventoryEvent) error {
    h.LogEvent(event, "handling stock alert event")
    
    data, ok := event.Data.(*types.StockAlertData)
    if !ok {
        return fmt.Errorf("invalid event data type")
    }

    if data.Available <= data.Threshold {
        req := &messageservice.SendTemplateMessageRequest{
            TemplateCode: "stock_alert",
            UserId:      event.Metadata.UserID,
            Params: map[string]string{
                "skuId":      fmt.Sprint(data.SkuID),
                "warehouse":  fmt.Sprint(data.WarehouseID),
                "available":  fmt.Sprint(data.Available),
                "threshold":  fmt.Sprint(data.Threshold),
            },
            Channels: []int32{1, 2}, // SMS and Email
        }

        _, err := h.messageRpc.SendTemplateMessage(context.Background(), req)
        if err != nil {
            h.LogError(event, err)
            return types.NewRetryableError(err)
        }
    }

    h.LogEvent(event, "stock alert event handled successfully")
    return nil
}