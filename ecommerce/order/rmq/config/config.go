package config

import "fmt"

type RabbitMQConfig struct {
	Host      string          `yaml:"Host"`
	Port      int             `yaml:"Port"`
	Username  string          `yaml:"Username"`
	Password  string          `yaml:"Password"`
	VHost     string          `yaml:"VHost"`
	Exchanges ExchangeConfigs `yaml:"Exchanges"`
	Queues    QueueConfigs    `yaml:"Queues"`

	// New configurations
	Retry struct {
		MaxAttempts     int     `yaml:"MaxAttempts"`
		InitialInterval int     `yaml:"InitialInterval"` // milliseconds
		MaxInterval     int     `yaml:"MaxInterval"`     // milliseconds
		BackoffFactor   float64 `yaml:"BackoffFactor"`
		Jitter          bool    `yaml:"Jitter"`
	} `yaml:"Retry"`

	Batch struct {
		Size          int `yaml:"Size"`
		FlushInterval int `yaml:"FlushInterval"` // milliseconds
		Workers       int `yaml:"Workers"`
	} `yaml:"Batch"`

	DeadLetter struct {
		Exchange   string `yaml:"Exchange"`
		Queue      string `yaml:"Queue"`
		RoutingKey string `yaml:"RoutingKey"`
	} `yaml:"DeadLetter"`

	Middleware struct {
		EnableRecovery bool `yaml:"EnableRecovery"`
		EnableLogging  bool `yaml:"EnableLogging"`
	} `yaml:"Middleware"`

	Server struct {
		Name      string           `yaml:"Name"`
		Mode      string           `yaml:"Mode"`
		LogLevel  string           `yaml:"LogLevel"`
		Consumers []ConsumerConfig `yaml:"Consumers"`
		Monitor   struct {
			Enabled bool `yaml:"Enabled"`
			Port    int  `yaml:"Port"`
		} `yaml:"Monitor"`
	} `yaml:"Server"`
}

type ConsumerConfig struct {
	Queue   string `yaml:"Queue"`
	Workers int    `yaml:"Workers"`
}

type ExchangeConfigs struct {
	OrderEvent ExchangeConfig `yaml:"OrderEvent"`
}

type ExchangeConfig struct {
	Name    string `yaml:"Name"`
	Type    string `yaml:"Type"`
	Durable bool   `yaml:"Durable"`
}

type QueueConfigs struct {
	OrderStatus QueueConfig `yaml:"OrderStatus"`
	OrderAlert  QueueConfig `yaml:"OrderAlert"`
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
