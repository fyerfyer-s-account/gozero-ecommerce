package config

import "strconv"

type RabbitMQConfig struct {
	Host      string          `yaml:"Host"`
	Port      int             `yaml:"Port"`
	Username  string          `yaml:"Username"`
	Password  string          `yaml:"Password"`
	VHost     string          `yaml:"VHost"`
	Exchanges ExchangeConfigs `yaml:"Exchanges"`
	Queues    QueueConfigs    `yaml:"Queues"`
}

type ExchangeConfigs struct {
	ProductEvent ExchangeConfig `yaml:"ProductEvent"`
}

type ExchangeConfig struct {
	Name    string `yaml:"Name"`
	Type    string `yaml:"Type"`
	Durable bool   `yaml:"Durable"`
}

type QueueConfigs struct {
	ProductUpdate QueueConfig `yaml:"ProductUpdate"`
	ProductStock  QueueConfig `yaml:"ProductStock"`
	ProductPrice  QueueConfig `yaml:"ProductPrice"`
}

type QueueConfig struct {
	Name       string `yaml:"Name"`
	RoutingKey string `yaml:"RoutingKey"`
	Durable    bool   `yaml:"Durable"`
}

func NewConfig() *RabbitMQConfig {
	return &RabbitMQConfig{
		Host:     "localhost",
		Port:     5672,
		Username: "guest",
		Password: "guest",
		VHost:    "/",
		Exchanges: ExchangeConfigs{
			ProductEvent: ExchangeConfig{
				Name:    "product.events",
				Type:    "topic",
				Durable: true,
			},
		},
		Queues: QueueConfigs{
			ProductUpdate: QueueConfig{
				Name:       "product.update",
				RoutingKey: "product.update.*",
				Durable:    true,
			},
			ProductStock: QueueConfig{
				Name:       "product.stock",
				RoutingKey: "product.stock.*",
				Durable:    true,
			},
			ProductPrice: QueueConfig{
				Name:       "product.price",
				RoutingKey: "product.price.*",
				Durable:    true,
			},
		},
	}
}

func (c *RabbitMQConfig) GetDSN() string {
	return "amqp://" + c.Username + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + c.VHost
}
