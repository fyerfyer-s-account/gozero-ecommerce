package config

import (
    "fmt"
    "time"
)

type RabbitMQConfig struct {
    Host     string `json:"host"`
    Port     int    `json:"port"`
    Username string `json:"username"`
    Password string `json:"password"`
    VHost    string `json:"vhost"`

    // Connection settings
    ConnectionTimeout time.Duration `json:"connection_timeout"`
    HeartbeatInterval time.Duration `json:"heartbeat_interval"`

    // Channel settings
    PrefetchCount  int  `json:"prefetch_count"`
    PrefetchGlobal bool `json:"prefetch_global"`

    // Exchanges configuration
    Exchanges []ExchangeConfig `json:"exchanges"`

    // Queues configuration
    Queues []QueueConfig `json:"queues"`
}

type ExchangeConfig struct {
    Name       string `json:"name"`
    Type       string `json:"type"`
    Durable    bool   `json:"durable"`
    AutoDelete bool   `json:"auto_delete"`
    Internal   bool   `json:"internal"`
    NoWait     bool   `json:"no_wait"`
}

type QueueConfig struct {
    Name       string            `json:"name"`
    Durable    bool             `json:"durable"`
    AutoDelete bool             `json:"auto_delete"`
    Exclusive  bool             `json:"exclusive"`
    NoWait     bool             `json:"no_wait"`
    Args       map[string]interface{} `json:"args"`
    Bindings   []BindingConfig  `json:"bindings"`
}

type BindingConfig struct {
    Exchange   string            `json:"exchange"`
    RoutingKey string           `json:"routing_key"`
    NoWait     bool             `json:"no_wait"`
    Args       map[string]interface{} `json:"args"`
}

func (c *RabbitMQConfig) DSN() string {
    return fmt.Sprintf("amqp://%s:%s@%s:%d/%s",
        c.Username,
        c.Password,
        c.Host,
        c.Port,
        c.VHost,
    )
}

func (c *RabbitMQConfig) Validate() error {
    if c.Host == "" {
        return fmt.Errorf("rabbitmq host is required")
    }
    if c.Port == 0 {
        c.Port = 5672 // default port
    }
    if c.ConnectionTimeout == 0 {
        c.ConnectionTimeout = 10 * time.Second
    }
    if c.HeartbeatInterval == 0 {
        c.HeartbeatInterval = 10 * time.Second
    }
    if c.PrefetchCount == 0 {
        c.PrefetchCount = 1
    }
    return nil
}

func DefaultConfig() *RabbitMQConfig {
    return &RabbitMQConfig{
        Host:              "localhost",
        Port:             5672,
        VHost:            "/",
        Username:         "guest",
        Password:         "guest",
        ConnectionTimeout: 10 * time.Second,
        HeartbeatInterval: 10 * time.Second,
        PrefetchCount:    1,
        PrefetchGlobal:   false,
        Exchanges:        make([]ExchangeConfig, 0),
        Queues:           make([]QueueConfig, 0),
    }
}