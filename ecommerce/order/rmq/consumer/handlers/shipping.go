package handlers

import (
    "context"
    "fmt"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/middleware"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
    "time"
)

type ShippingHandler struct {
    BaseHandler
    messageRpc         messageservice.MessageService
    orderModel         model.OrdersModel
    orderShippingModel model.OrderShippingModel
}

func NewShippingHandler(
    logger middleware.Logger,
    messageRpc messageservice.MessageService,
    orderModel model.OrdersModel,
    orderShippingModel model.OrderShippingModel,
) *ShippingHandler {
    return &ShippingHandler{
        BaseHandler:        NewBaseHandler(logger),
        messageRpc:         messageRpc,
        orderModel:         orderModel,
        orderShippingModel: orderShippingModel,
    }
}

func (h *ShippingHandler) Handle(event *types.OrderEvent) error {
    h.LogEvent(event, "handling shipping event")

    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    var err error
    switch event.Type {
    case types.EventTypeOrderShipped:
        err = h.handleOrderShipped(ctx, event)
    default:
        h.LogEvent(event, "ignoring irrelevant event type")
        return nil
    }

    if err != nil {
        h.LogError(event, err)
        return types.NewRetryableError(err)
    }

    h.LogEvent(event, "shipping event handled successfully")
    return nil
}

func (h *ShippingHandler) handleOrderShipped(ctx context.Context, event *types.OrderEvent) error {
    data, ok := event.Data.(*types.OrderShippedData)
    if !ok {
        return fmt.Errorf("invalid event data type: %T", event.Data)
    }

    h.LogEvent(event, "processing order shipped event",
        "order_id", data.OrderId,
        "shipping_no", data.ShippingNo,
        "company", data.Company,
    )

    // 1. Find order shipping record
    shipping, err := h.orderShippingModel.FindByOrderId(ctx, uint64(data.OrderId))
    if err != nil {
        h.LogError(event, fmt.Errorf("failed to find shipping record: %w", err))
        return err
    }

    // 2. Update shipping info
    if err := h.orderShippingModel.UpdateShippingInfo(ctx,
        shipping.OrderId, data.ShippingNo, data.Company); err != nil {
        h.LogError(event, fmt.Errorf("failed to update shipping info: %w", err))
        return err
    }

    // 3. Send notification with retry
    if err := h.sendShippingNotification(ctx, event, data); err != nil {
        h.LogError(event, fmt.Errorf("failed to send notification: %w", err))
        return err
    }

    h.LogEvent(event, "order shipped event processed successfully",
        "order_id", data.OrderId,
    )
    return nil
}

func (h *ShippingHandler) sendShippingNotification(ctx context.Context, event *types.OrderEvent, data *types.OrderShippedData) error {
    _, err := h.messageRpc.SendTemplateMessage(ctx, &messageservice.SendTemplateMessageRequest{
        TemplateCode: "order_shipped",
        UserId:       event.Metadata.UserID,
        Params: map[string]string{
            "orderNo":    data.OrderNo,
            "shippingNo": data.ShippingNo,
            "company":    data.Company,
        },
    })

    if err != nil {
        return fmt.Errorf("failed to send shipping notification: %w", err)
    }

    return nil
}