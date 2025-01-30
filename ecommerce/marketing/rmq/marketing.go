package main

import (
    "context"
    "flag"
    "os"
    "os/signal"
    "syscall"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/internal/server"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/internal/svc"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/zeromicro/go-zero/core/conf"
    "github.com/zeromicro/go-zero/core/service"
)

var (
    configFile = flag.String("f", "etc/marketing.yaml", "config file path")
    logger     = zerolog.GetLogger()
)

func main() {
    flag.Parse()

    // Create root context with cancellation
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Load configuration
    logger.Info(ctx, "Loading configuration...", map[string]interface{}{
        "configFile": *configFile,
    })

    var c config.Config
    conf.MustLoad(*configFile, &c)

    // Initialize service context
    logger.Info(ctx, "Initializing service context...", nil)
    serviceContext := svc.NewServiceContext(c)

    // Initialize RMQ server
    logger.Info(ctx, "Initializing RMQ server...", nil)
    rmqServer := server.NewRmqServer(ctx, serviceContext)

    // Create service group
    group := service.NewServiceGroup()
    group.Add(rmqServer)

    // Setup signal handling for graceful shutdown
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

    // Start service group
    logger.Info(ctx, "Starting marketing RMQ service...", map[string]interface{}{
        "listenOn": c.ListenOn,
    })

    go func() {
        group.Start()
    }()

    // Wait for shutdown signal
    select {
    case sig := <-sigChan:
        logger.Info(ctx, "Received shutdown signal", map[string]interface{}{
            "signal": sig.String(),
        })
    case <-ctx.Done():
        logger.Info(ctx, "Context cancelled", nil)
    }

    // Graceful shutdown
    logger.Info(ctx, "Initiating graceful shutdown...", nil)
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()

    go func() {
        group.Stop()
        shutdownCancel()
    }()

    // Wait for graceful shutdown or timeout
    <-shutdownCtx.Done()
    if shutdownCtx.Err() == context.DeadlineExceeded {
        logger.Warn(ctx, "Graceful shutdown timed out", nil)
    } else {
        logger.Info(ctx, "Graceful shutdown completed", nil)
    }
}