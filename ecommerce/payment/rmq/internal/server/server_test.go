package server

import (
    "context"
    "testing"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rmq/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/payment/rmq/internal/svc"
    "github.com/stretchr/testify/assert"
    "github.com/zeromicro/go-zero/core/conf"
)

func TestDatabaseConnection(t *testing.T) {
    // Load config
    var c config.Config
    conf.MustLoad("../../etc/payment.yaml", &c)
    
    t.Logf("Testing connection to: %s", c.Mysql.DataSource)

    // Create service context
    svcCtx := svc.NewServiceContext(c)
    assert.NotNil(t, svcCtx)

    // Try a simple query to verify connection
    ctx := context.Background()
    _, err := svcCtx.PaymentOrdersModel.FindOne(ctx, 1)
    if err != nil && err.Error() != "sql: no rows in result set" {
        // We expect either a valid result or "no rows" error
        // Any other error indicates connection issues 
        t.Errorf("Database connection failed: %v", err)
    }
}

func TestRabbitMQConnection(t *testing.T) {
    // Load config  
    var c config.Config
    conf.MustLoad("../../etc/payment.yaml", &c)

    // Create service context
    svcCtx := svc.NewServiceContext(c)
    assert.NotNil(t, svcCtx)

    // Check RabbitMQ connection
    assert.True(t, svcCtx.Broker.IsConnected())

    // Check channel
    assert.NotNil(t, svcCtx.Channel)
}