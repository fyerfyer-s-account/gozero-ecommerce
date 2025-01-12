package handlers

import (
	"context"
	"fmt"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
)

type StockUpdateHandler struct {
    inventoryRpc inventoryclient.Inventory
}

func NewStockUpdateHandler(inventoryRpc inventoryclient.Inventory) *StockUpdateHandler {
    return &StockUpdateHandler{
        inventoryRpc: inventoryRpc,
    }
}

func (h *StockUpdateHandler) Handle(event *types.InventoryEvent) error {
    data, ok := event.Data.(*types.StockUpdateData)
    if !ok {
        return fmt.Errorf("invalid event data type")
    }

    _, err := h.inventoryRpc.UpdateStock(context.Background(), &inventoryclient.UpdateStockRequest{
        SkuId:       int64(data.SkuID),
        WarehouseId: int64(data.WarehouseID),
        Quantity:    data.Quantity,
        Remark:      data.Remark,
    })

    return err
}