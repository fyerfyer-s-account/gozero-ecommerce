package config

import "fmt"

type RabbitMQConfig struct {
    Host      string          `yaml:"Host"`
    Port      int            `yaml:"Port"`
    Username  string         `yaml:"Username"`
    Password  string         `yaml:"Password"`
    VHost     string         `yaml:"VHost"`
    Exchanges ExchangeConfigs `yaml:"Exchanges"`
    Queues    QueueConfigs    `yaml:"Queues"`

    Retry struct {
        MaxAttempts     int     `yaml:"MaxAttempts"`
        InitialInterval int     `yaml:"InitialInterval"` // milliseconds
        MaxInterval     int     `yaml:"MaxInterval"`     // milliseconds
        BackoffFactor   float64 `yaml:"BackoffFactor"`
        Jitter         bool    `yaml:"Jitter"`
    } `yaml:"Retry"`

    Batch struct {
        Size          int `yaml:"Size"`
        FlushInterval int `yaml:"FlushInterval"` // milliseconds
        Workers       int `yaml:"Workers"`
    } `yaml:"Batch"`

    DeadLetter struct {
        Exchange     string `yaml:"Exchange"`
        Queue       string `yaml:"Queue"`
        RoutingKey  string `yaml:"RoutingKey"`
    } `yaml:"DeadLetter"`

    Middleware struct {
        EnableRecovery bool `yaml:"EnableRecovery"`
        EnableLogging  bool `yaml:"EnableLogging"`
    } `yaml:"Middleware"`
}

type ExchangeConfigs struct {
    InventoryEvent ExchangeConfig `yaml:"InventoryEvent"`
}

type ExchangeConfig struct {
    Name    string `yaml:"Name"`
    Type    string `yaml:"Type"`
    Durable bool   `yaml:"Durable"`
}

type QueueConfigs struct {
    StockUpdate QueueConfig `yaml:"StockUpdate"`
    StockAlert  QueueConfig `yaml:"StockAlert"`
    StockLock   QueueConfig `yaml:"StockLock"`
    OrderEvents QueueConfig `yaml:"OrderEvents"` // New: Order events queue
}

type QueueConfig struct {
    Name       string `yaml:"Name"`
    RoutingKey string `yaml:"RoutingKey"`
    Durable    bool   `yaml:"Durable"`
}

func (c *RabbitMQConfig) GetDSN() string {
    return fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
        c.Username,
        c.Password,
        c.Host,
        c.Port,
        c.VHost,
    )
}