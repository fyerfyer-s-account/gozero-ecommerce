package setup

import (
	"context"
	"testing"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/order-payment/helpers"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zeromicro/go-zero/core/conf"
)

func TestConfigLoading(t *testing.T) {
    var c Config
    err := conf.Load("../config/test.yaml", &c)
    require.NoError(t, err)

    // Verify essential config fields
    assert.NotEmpty(t, c.Mysql.DataSource)
    assert.NotEmpty(t, c.RabbitMQ.Host)
    assert.NotEmpty(t, c.CacheRedis)
}

func TestDatabaseConnection(t *testing.T) {
    var c Config
    conf.MustLoad("../config/test.yaml", &c)

    // Create DB helper
    db, err := helpers.NewDBHelper(c.Mysql.DataSource, c.CacheRedis)
    require.NoError(t, err)
    defer db.Close()

    // Try a simple query
    ctx := context.Background()
    _, err = db.GetOrdersModel().FindOne(ctx, 1)
    if err != nil && err.Error() != "sql: no rows in result set" {
        t.Errorf("Database connection failed: %v", err)
    }
}

func TestRabbitMQConnection(t *testing.T) {
    var c Config
    conf.MustLoad("../config/test.yaml", &c)

    // Create RMQ config
    rmqConfig := &config.RabbitMQConfig{
        Host:              c.RabbitMQ.Host,
        Port:              c.RabbitMQ.Port,
        Username:          c.RabbitMQ.Username,
        Password:          c.RabbitMQ.Password,
        VHost:             c.RabbitMQ.VHost,
        ConnectionTimeout: time.Duration(c.RabbitMQ.ConnectionTimeout),
        HeartbeatInterval: time.Duration(c.RabbitMQ.HeartbeatInterval),
    }

    // Create RMQ helper
    rmq, err := helpers.NewRMQHelper(rmqConfig)
    require.NoError(t, err)
    defer rmq.Close()

    // Verify connection
    assert.True(t, rmq.IsConnected())
    assert.NotNil(t, rmq.GetChannel())
}

func TestNewTestContext(t *testing.T) {
    ctx, err := NewTestContext()
    require.NoError(t, err)
    defer ctx.Close()

    // Verify context components
    assert.NotNil(t, ctx.Config)
    assert.NotNil(t, ctx.DB)
    assert.NotNil(t, ctx.RMQ)

    // Verify RMQ connection
    assert.True(t, ctx.RMQ.IsConnected())

    // Try database query
    dbCtx := context.Background()
    _, err = ctx.DB.GetOrdersModel().FindOne(dbCtx, 1)
    if err != nil && err.Error() != "sql: no rows in result set" {
        t.Errorf("Database connection failed: %v", err)
    }

    // Try RMQ operations
    err = ctx.RMQ.PublishEvent("order.events", "test.event", map[string]string{
        "test": "data",
    })
    require.NoError(t, err)
}