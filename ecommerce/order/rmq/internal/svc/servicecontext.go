package svc

import (
	"context"
	"time"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/consumer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rmq/producer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/order/rpc/model"
	rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/broker"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
	Config config.Config
	Broker broker.Broker

	// Models
	OrdersModel        model.OrdersModel
	OrderItemsModel    model.OrderItemsModel
	OrderShippingModel model.OrderShippingModel
	OrderRefundsModel  model.OrderRefundsModel
	OrderPaymentsModel model.OrderPaymentsModel

	// RMQ Components
	Producer *producer.OrderProducer
	Consumer *consumer.OrderConsumer
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	// Initialize RabbitMQ broker
	rmqConfig := convertToEventbusConfig(&c)
	rmqBroker := broker.NewAMQPBroker(rmqConfig)
	if err := rmqBroker.Connect(context.Background()); err != nil {
		panic(err)
	}

	// Get RMQ channel
	ch, err := rmqBroker.Channel()
	if err != nil {
		panic(err)
	}

	// Initialize models
	ordersModel := model.NewOrdersModel(conn, c.CacheRedis)
	orderItemsModel := model.NewOrderItemsModel(conn, c.CacheRedis)
	orderShippingModel := model.NewOrderShippingModel(conn, c.CacheRedis)
	orderRefundsModel := model.NewOrderRefundsModel(conn, c.CacheRedis)
	orderPaymentsModel := model.NewOrderPaymentsModel(conn, c.CacheRedis)

	// Initialize producer and consumer
	prod := producer.NewOrderProducer(ch, "order.events")
	cons := consumer.NewOrderConsumer(
		ch,
		ordersModel,
		orderPaymentsModel,
		orderShippingModel,
		orderRefundsModel,
	)

	return &ServiceContext{
		Config:             c,
		Broker:             rmqBroker,
		OrdersModel:        ordersModel,
		OrderItemsModel:    orderItemsModel,
		OrderShippingModel: orderShippingModel,
		OrderRefundsModel:  orderRefundsModel,
		OrderPaymentsModel: orderPaymentsModel,
		Producer:           prod,
		Consumer:           cons,
	}
}

func convertToEventbusConfig(c *config.Config) *rmqconfig.RabbitMQConfig {
	exchanges := make([]rmqconfig.ExchangeConfig, len(c.RabbitMQ.Exchanges))
	for i, e := range c.RabbitMQ.Exchanges {
		exchanges[i] = rmqconfig.ExchangeConfig{
			Name:       e.Name,
			Type:       e.Type,
			Durable:    e.Durable,
			AutoDelete: e.AutoDelete,
			Internal:   e.Internal,
			NoWait:     e.NoWait,
		}
	}

	queues := make([]rmqconfig.QueueConfig, len(c.RabbitMQ.Queues))
	for i, q := range c.RabbitMQ.Queues {
		bindings := make([]rmqconfig.BindingConfig, len(q.Bindings))
		for j, b := range q.Bindings {
			bindings[j] = rmqconfig.BindingConfig{
				Exchange:   b.Exchange,
				RoutingKey: b.RoutingKey,
				NoWait:     b.NoWait,
			}
		}
		queues[i] = rmqconfig.QueueConfig{
			Name:       q.Name,
			Durable:    q.Durable,
			AutoDelete: q.AutoDelete,
			Exclusive:  q.Exclusive,
			NoWait:     q.NoWait,
			Bindings:   bindings,
		}
	}

	return &rmqconfig.RabbitMQConfig{
		Host:              c.RabbitMQ.Host,
		Port:              c.RabbitMQ.Port,
		Username:          c.RabbitMQ.Username,
		Password:          c.RabbitMQ.Password,
		VHost:             c.RabbitMQ.VHost,
		ConnectionTimeout: time.Duration(c.RabbitMQ.ConnectionTimeout) * time.Second,
		HeartbeatInterval: time.Duration(c.RabbitMQ.HeartbeatInterval) * time.Second,
		PrefetchCount:     c.RabbitMQ.PrefetchCount,
		PrefetchGlobal:    c.RabbitMQ.PrefetchGlobal,
		Exchanges:         exchanges,
		Queues:            queues,
	}
}
