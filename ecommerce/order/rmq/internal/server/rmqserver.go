package server

import (
    "context"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
)

type RmqServer struct {
    ctx    context.Context
    svcCtx *svc.ServiceContext
    logger *zerolog.Logger
}

func NewRmqServer(ctx context.Context, svcCtx *svc.ServiceContext) *RmqServer {
    return &RmqServer{
        ctx:    ctx,
        svcCtx: svcCtx,
        logger: zerolog.GetLogger(),
    }
}

// Start implements service.Service
func (s *RmqServer) Start() {
    // Start the consumer
    if err := s.svcCtx.Consumer.Start(s.ctx); err != nil {
        s.logger.Error(s.ctx, "Failed to start consumer", err, nil)
        panic(err) // In production you might want to handle this differently
    }

    s.logger.Info(s.ctx, "Order RMQ server started", nil)
}

// Stop implements service.Service
func (s *RmqServer) Stop() {
    // Stop the consumer
    if err := s.svcCtx.Consumer.Stop(); err != nil {
        s.logger.Error(s.ctx, "Failed to stop consumer", err, nil)
    }

    // Close the producer
    if err := s.svcCtx.Producer.Close(); err != nil {
        s.logger.Error(s.ctx, "Failed to close producer", err, nil)
    }

    // Disconnect broker
    if err := s.svcCtx.Broker.Disconnect(); err != nil {
        s.logger.Error(s.ctx, "Failed to disconnect broker", err, nil)
    }

    s.logger.Info(s.ctx, "Order RMQ server stopped", nil)
}