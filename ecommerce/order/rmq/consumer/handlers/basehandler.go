package handlers

import (
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/types"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/middleware"
)

type BaseHandler struct {
    logger middleware.Logger
}

func NewBaseHandler(logger middleware.Logger) BaseHandler {
    return BaseHandler{logger: logger}
}

func (h *BaseHandler) LogEvent(event *types.OrderEvent, msg string, args ...interface{}) {
    h.logger.Info(msg,
        "event_id", event.ID,
        "event_type", event.Type,
        "trace_id", event.Metadata.TraceID,
        "args", args,
    )
}

func (h *BaseHandler) LogError(event *types.OrderEvent, err error) {
    h.logger.Error("handler error",
        "event_id", event.ID,
        "event_type", event.Type,
        "trace_id", event.Metadata.TraceID,
        "error", err,
    )
}