package svc

import (
	"context"

	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/consumer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/producer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/inventoryclient"
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

    // Initialize RPC clients
    messageRpc := messageservice.NewMessageService(zrpc.MustNewClient(c.MessageRpc))

    // Initialize producer
    prod, err := producer.NewProducer(&c.RabbitMQ)
    if err != nil {
        panic(err)
    }

	inventoryRpc := inventoryclient.NewInventory(zrpc.MustNewClient(c.Etcd.RpcClientConf))
	logger := &logWrapper{Logger: logx.WithContext(context.TODO())}
    // Initialize consumer with logger and RPC clients
    cons, err := consumer.NewConsumer(
		&c.RabbitMQ,
		logger,
		inventoryRpc,
		messageRpc,
	)
	if err != nil {
		panic(err)
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
    if err := cons.Start(); err != nil {
        panic(err)
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