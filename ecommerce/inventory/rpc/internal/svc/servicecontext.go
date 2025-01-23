package svc

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/consumer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/consumer/handlers"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/producer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/types"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config            config.Config
	StocksModel       model.StocksModel
	StockLocksModel   model.StockLocksModel
	StockRecordsModel model.StockRecordsModel
	WarehousesModel   model.WarehousesModel
	Producer          *producer.Producer
	Consumer          *consumer.Consumer
	MessageRpc        messageservice.MessageService
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)
	logger := &logWrapper{Logger: logx.WithContext(context.TODO())}

	// Initialize RPC clients
	messageRpc := messageservice.NewMessageService(zrpc.MustNewClient(c.MessageRpc))

	// Initialize producer
	prod, err := producer.NewProducer(&c.RabbitMQ)
	if err != nil {
		panic(err)
	}

	// Initialize consumer if in RMQ server mode
	var cons *consumer.Consumer
	if c.RabbitMQ.Server.Mode != "" {
		cons, err = consumer.NewConsumer(&c.RabbitMQ, logger, nil, messageRpc)
		if err != nil {
			panic(err)
		}

		// Register handlers only in RMQ server mode
		orderHandler := handlers.NewOrderEventHandler(logger, nil)
		cons.Subscribe(types.EventTypeOrderCreated, orderHandler)
		cons.Subscribe(types.EventTypeOrderCancelled, orderHandler)
		cons.Subscribe(types.EventTypeOrderPaid, orderHandler)
		cons.Subscribe(types.EventTypeOrderRefunded, orderHandler)
	}

	svcCtx := &ServiceContext{
		Config:            c,
		StocksModel:       model.NewStocksModel(conn, c.CacheRedis),
		StockLocksModel:   model.NewStockLocksModel(conn, c.CacheRedis),
		StockRecordsModel: model.NewStockRecordsModel(conn, c.CacheRedis),
		WarehousesModel:   model.NewWarehousesModel(conn, c.CacheRedis),
		Producer:          prod,
		Consumer:          cons,
		MessageRpc:        messageRpc,
	}

	// Start consumer
	if cons != nil {
		if err := cons.Start(); err != nil {
			panic(err)
		}
	}

	return svcCtx
}

type logWrapper struct {
	logx.Logger
}

func (l *logWrapper) Info(msg string, keysAndValues ...interface{}) {
	l.Logger.Infof(msg, keysAndValues...)
}

func (l *logWrapper) Error(msg string, keysAndValues ...interface{}) {
	l.Logger.Errorf(msg, keysAndValues...)
}
