package svc

import (
    "context"
    "fmt"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/broker"
    rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
    Config config.Config

    MessagesModel      model.MessagesModel
    MessageSendsModel  model.MessageSendsModel
    MessageTemplatesModel     model.MessageTemplatesModel
    SettingsModel      model.NotificationSettingsModel

    Producer *producer.MessageProducer
}

func NewServiceContext(c config.Config) *ServiceContext {
    // Initialize MySQL connection
    conn := sqlx.NewMysql(c.Mysql.DataSource)

    // Initialize RabbitMQ broker
    rmqBroker := broker.NewAMQPBroker(&rmqconfig.RabbitMQConfig{
        Host:              c.RabbitMQ.Host,
        Port:              c.RabbitMQ.Port,
        Username:          c.RabbitMQ.Username,
        Password:          c.RabbitMQ.Password,
        VHost:             c.RabbitMQ.VHost,
        ConnectionTimeout: time.Duration(c.RabbitMQ.ConnectionTimeout) * time.Second,
        HeartbeatInterval: time.Duration(c.RabbitMQ.HeartbeatInterval) * time.Second,
    })

    // Establish RabbitMQ connection
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    if err := rmqBroker.Connect(ctx); err != nil {
        panic(fmt.Sprintf("Failed to connect to RabbitMQ broker: %v", err))
    }

    // Get RabbitMQ channel
    ch, err := rmqBroker.Channel()
    if err != nil {
        panic(fmt.Sprintf("Failed to create RabbitMQ channel: %v", err))
    }

    return &ServiceContext{
        Config: c,

        MessagesModel:      model.NewMessagesModel(conn, c.CacheRedis),
        MessageSendsModel:  model.NewMessageSendsModel(conn, c.CacheRedis),
        MessageTemplatesModel:     model.NewMessageTemplatesModel(conn, c.CacheRedis),
        SettingsModel:      model.NewNotificationSettingsModel(conn, c.CacheRedis),

        Producer: producer.NewMessageProducer(ch, c.RabbitMQ.Exchange),
    }
}