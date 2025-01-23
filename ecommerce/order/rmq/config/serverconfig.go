package config

import (
    "github.com/zeromicro/go-zero/core/service"
)

type Config struct {
    service.ServiceConf // Includes Name and Mode
    RabbitMQ RabbitMQConfig `yaml:"RabbitMQ"`
}