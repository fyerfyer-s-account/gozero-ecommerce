package setup

import (
    "os"
    "time"
)

// TestConfig holds all test configuration
type TestConfig struct {
    RabbitMQURL      string
    TestTimeout      time.Duration
    MessageWaitTime  time.Duration
}

// NewTestConfig creates a new test configuration with default values
func NewTestConfig() *TestConfig {
    return &TestConfig{
        RabbitMQURL:      getEnvOrDefault("TEST_RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
        TestTimeout:      getEnvDurationOrDefault("TEST_TIMEOUT", 30*time.Second),
        MessageWaitTime:  getEnvDurationOrDefault("TEST_MESSAGE_WAIT", 5*time.Second),
    }
}

// getEnvOrDefault returns environment variable value or default if not set
func getEnvOrDefault(key, defaultValue string) string {
    if value, exists := os.LookupEnv(key); exists {
        return value
    }
    return defaultValue
}

// getEnvDurationOrDefault returns duration from environment variable or default if not set
func getEnvDurationOrDefault(key string, defaultValue time.Duration) time.Duration {
    if value, exists := os.LookupEnv(key); exists {
        if duration, err := time.ParseDuration(value); err == nil {
            return duration
        }
    }
    return defaultValue
}

// TestContext represents the context for a test case
type TestContext struct {
    Config  *TestConfig
    RMQ     *RMQTestSetup
}

// NewTestContext creates a new test context
func NewTestContext() (*TestContext, error) {
    config := NewTestConfig()
    
    rmq, err := NewRMQTestSetup(config)
    if err != nil {
        return nil, err
    }

    return &TestContext{
        Config: config,
        RMQ:    rmq,
    }, nil
}

// Close cleans up test context resources
func (tc *TestContext) Close() error {
    if tc.RMQ != nil {
        if err := tc.RMQ.Close(); err != nil {
            return err
        }
    }
    return nil
}

// CleanupTestData cleans up any test data
func (tc *TestContext) CleanupTestData() error {
    if tc.RMQ != nil {
        if err := tc.RMQ.CleanupTestQueues(); err != nil {
            return err
        }
    }
    return nil
}