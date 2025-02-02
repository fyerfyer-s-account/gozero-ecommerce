package setup

import (
	"flag"
	"fmt"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/e2e/rmq/order-payment/helpers"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"	
	"github.com/zeromicro/go-zero/core/conf"
)

var (
    configFile = flag.String("f", "../config/test.yaml", "config file path")
    // logger     = zerolog.GetLogger()
)

// TestContext holds all the dependencies required for e2e testing
type TestContext struct {
    Config *Config
    DB     *helpers.DBHelper
    RMQ    *helpers.RMQHelper
}

// NewTestContext creates and initializes a new test context
func NewTestContext() (*TestContext, error) {
    // Load configuration
    var c Config
    if err := conf.Load(*configFile, &c); err != nil {
        return nil, fmt.Errorf("failed to load config: %w", err)
    }

    // Initialize DB helper
    db, err := helpers.NewDBHelper(c.Mysql.DataSource, c.CacheRedis)
    if err != nil {
        return nil, fmt.Errorf("failed to create DB helper: %w", err)
    }

    // Initialize RMQ helper
    rmqConfig := &config.RabbitMQConfig{
        Host:              c.RabbitMQ.Host,
        Port:              c.RabbitMQ.Port,
        Username:          c.RabbitMQ.Username,
        Password:          c.RabbitMQ.Password,
        VHost:             c.RabbitMQ.VHost,
        ConnectionTimeout: time.Duration(c.RabbitMQ.ConnectionTimeout),
        HeartbeatInterval: time.Duration(c.RabbitMQ.HeartbeatInterval),
    }
    
    rmq, err := helpers.NewRMQHelper(rmqConfig)
    if err != nil {
        db.Close() // Clean up DB connection if RMQ fails
        return nil, fmt.Errorf("failed to create RMQ helper: %w", err)
    }

    return &TestContext{
        Config: &c,
        DB:     db,
        RMQ:    rmq,
    }, nil
}

// Close releases all resources held by the test context
func (tc *TestContext) Close() error {
    var errs []error

    // Clean up any test queues first
    if tc.RMQ != nil && tc.RMQ.GetChannel() != nil {
        // Channel is already closed in RMQ.Close()
    }

    // Close RMQ connections
    if tc.RMQ != nil {
        if err := tc.RMQ.Close(); err != nil {
            errs = append(errs, fmt.Errorf("failed to close RMQ: %w", err))
        }
    }

    // Close DB connections
    if tc.DB != nil {
        if err := tc.DB.Close(); err != nil {
            errs = append(errs, fmt.Errorf("failed to close DB: %w", err))
        }
    }

    if len(errs) > 0 {
        return fmt.Errorf("errors closing test context: %v", errs)
    }
    return nil
}