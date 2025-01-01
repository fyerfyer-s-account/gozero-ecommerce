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
	UserEvent ExchangeConfig `yaml:"UserEvent"`
}

type ExchangeConfig struct {
	Name    string `yaml:"Name"`
	Type    string `yaml:"Type"`
	Durable bool   `yaml:"Durable"`
}

type QueueConfigs struct {
	UserNotification QueueConfig `yaml:"UserNotification"`
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
			UserEvent: ExchangeConfig{
				Name:    "user.events",
				Type:    "topic",
				Durable: true,
			},
		},
		Queues: QueueConfigs{
			UserNotification: QueueConfig{
				Name:       "user.notification",
				RoutingKey: "user.#",
				Durable:    true,
			},
		},
	}
}

func (c *RabbitMQConfig) GetDSN() string {
	return "amqp://" + c.Username + ":" + c.Password + "@" + c.Host + ":" + strconv.Itoa(c.Port) + "/" + c.VHost
}
