package svc

import (
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
    rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/consumer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/producer"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
    Config              config.Config
    MessagesModel       model.MessagesModel
    MessageSendsModel   model.MessageSendsModel
    MessageTemplatesModel model.MessageTemplatesModel
    NotificationSettingsModel model.NotificationSettingsModel
    Producer            *producer.Producer
    Consumer            *consumer.Consumer
}

func NewServiceContext(c config.Config) *ServiceContext {
    conn := sqlx.NewMysql(c.Mysql.DataSource)

    // Initialize RabbitMQ config
    rmqConfig := &rmqconfig.RabbitMQConfig{
        Host:     c.RabbitMQ.Host,
        Port:     c.RabbitMQ.Port,
        Username: c.RabbitMQ.Username,
        Password: c.RabbitMQ.Password,
        VHost:    c.RabbitMQ.VHost,
        Exchanges: rmqconfig.ExchangeConfigs{
            MessageEvent: rmqconfig.ExchangeConfig{
                Name:    c.RabbitMQ.Exchanges.MessageEvent.Name,
                Type:    c.RabbitMQ.Exchanges.MessageEvent.Type,
                Durable: c.RabbitMQ.Exchanges.MessageEvent.Durable,
            },
        },
        Queues: rmqconfig.QueueConfigs{
            NotificationQueue: rmqconfig.QueueConfig{
                Name:       c.RabbitMQ.Queues.NotificationQueue.Name,
                RoutingKey: c.RabbitMQ.Queues.NotificationQueue.RoutingKey,
                Durable:    c.RabbitMQ.Queues.NotificationQueue.Durable,
            },
            TemplateQueue: rmqconfig.QueueConfig{
                Name:       c.RabbitMQ.Queues.TemplateQueue.Name,
                RoutingKey: c.RabbitMQ.Queues.TemplateQueue.RoutingKey,
                Durable:    c.RabbitMQ.Queues.TemplateQueue.Durable,
            },
        },
    }

    // Initialize producer
    prod, err := producer.NewProducer(*rmqConfig)
    if err != nil {
        panic(err)
    }

    // Initialize consumer
    cons := consumer.NewConsumer(*rmqConfig)

    return &ServiceContext{
        Config:              c,
        MessagesModel:       model.NewMessagesModel(conn, c.CacheRedis),
        MessageSendsModel:   model.NewMessageSendsModel(conn, c.CacheRedis),
        MessageTemplatesModel: model.NewMessageTemplatesModel(conn, c.CacheRedis),
        NotificationSettingsModel: model.NewNotificationSettingsModel(conn, c.CacheRedis),
        Producer:            prod,
        Consumer:            cons,
    }
}