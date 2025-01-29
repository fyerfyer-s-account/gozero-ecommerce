package server

import (
    "context"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rmq/internal/svc"
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

func (s *RmqServer) Start() {
    s.logger.Info(s.ctx, "Starting Payment RMQ server...", nil)

    if !s.svcCtx.Broker.IsConnected() {
        if err := s.svcCtx.Broker.Connect(s.ctx); err != nil {
            s.logger.Error(s.ctx, "Failed to connect to RabbitMQ broker", err, nil)
            panic(err)
        }
    }

    if err := s.svcCtx.Consumer.Start(s.ctx); err != nil {
        s.logger.Error(s.ctx, "Failed to start consumer", err, nil)
        panic(err)
    }

    s.logger.Info(s.ctx, "Payment RMQ server started successfully", nil)
}

func (s *RmqServer) Stop() {
    s.logger.Info(s.ctx, "Stopping Payment RMQ server...", nil)

    if err := s.svcCtx.Consumer.Stop(); err != nil {
        s.logger.Error(s.ctx, "Failed to stop consumer", err, nil)
    }

    if err := s.svcCtx.Producer.Close(); err != nil {
        s.logger.Error(s.ctx, "Failed to close producer", err, nil)
    }

    if err := s.svcCtx.Broker.Close(); err != nil {
        s.logger.Error(s.ctx, "Failed to close broker", err, nil)
    }

    s.logger.Info(s.ctx, "Payment RMQ server stopped successfully", nil)
}