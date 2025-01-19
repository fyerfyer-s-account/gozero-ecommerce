package svc

import (
	rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/consumer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rmq/producer"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/internal/config"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/marketing/rpc/model"
	"github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/messageservice"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"
)

type ServiceContext struct {
	Config             config.Config
	CouponsModel       model.CouponsModel
	UserCouponsModel   model.UserCouponsModel
	PromotionsModel    model.PromotionsModel
	UserPointsModel    model.UserPointsModel
	PointsRecordsModel model.PointsRecordsModel
	Producer           *producer.Producer
	Consumer           *consumer.Consumer
	MessageRpc         messageservice.MessageService
}

func NewServiceContext(c config.Config) *ServiceContext {
	conn := sqlx.NewMysql(c.Mysql.DataSource)

	// Initialize RPC clients
	messageRpc := messageservice.NewMessageService(zrpc.MustNewClient(c.MessageRpc))

	// Initialize RabbitMQ config
	rmqConfig := &rmqconfig.RabbitMQConfig{
		Host:     c.RabbitMQ.Host,
		Port:     c.RabbitMQ.Port,
		Username: c.RabbitMQ.Username,
		Password: c.RabbitMQ.Password,
		VHost:    c.RabbitMQ.VHost,
		Exchanges: rmqconfig.ExchangeConfigs{
			MarketingEvent: rmqconfig.ExchangeConfig{
				Name:    c.RabbitMQ.Exchanges.MarketingEvent.Name,
				Type:    c.RabbitMQ.Exchanges.MarketingEvent.Type,
				Durable: c.RabbitMQ.Exchanges.MarketingEvent.Durable,
			},
		},
		Queues: rmqconfig.QueueConfigs{
			CouponEvent: rmqconfig.QueueConfig{
				Name:       c.RabbitMQ.Queues.CouponEvent.Name,
				RoutingKey: c.RabbitMQ.Queues.CouponEvent.RoutingKey,
				Durable:    c.RabbitMQ.Queues.CouponEvent.Durable,
			},
			PromotionEvent: rmqconfig.QueueConfig{
				Name:       c.RabbitMQ.Queues.PromotionEvent.Name,
				RoutingKey: c.RabbitMQ.Queues.PromotionEvent.RoutingKey,
				Durable:    c.RabbitMQ.Queues.PromotionEvent.Durable,
			},
			PointsEvent: rmqconfig.QueueConfig{
				Name:       c.RabbitMQ.Queues.PointsEvent.Name,
				RoutingKey: c.RabbitMQ.Queues.PointsEvent.RoutingKey,
				Durable:    c.RabbitMQ.Queues.PointsEvent.Durable,
			},
		},
	}

	// Initialize models
	couponsModel := model.NewCouponsModel(conn, c.CacheRedis)
	userCouponsModel := model.NewUserCouponsModel(conn, c.CacheRedis)
	promotionsModel := model.NewPromotionsModel(conn, c.CacheRedis)
	userPointsModel := model.NewUserPointsModel(conn, c.CacheRedis)
	pointsRecordsModel := model.NewPointsRecordsModel(conn, c.CacheRedis)

	// Initialize producer
	prod, err := producer.NewProducer(rmqConfig)
	if err != nil {
		panic(err)
	}

	// Initialize consumer
	cons, err := consumer.NewConsumer(
		rmqConfig,
		couponsModel,
		userCouponsModel,
		promotionsModel,
		userPointsModel,
		pointsRecordsModel,
	)
	if err != nil {
		panic(err)
	}

	svcCtx := &ServiceContext{
		Config:             c,
		CouponsModel:       couponsModel,
		UserCouponsModel:   userCouponsModel,
		PromotionsModel:    promotionsModel,
		UserPointsModel:    userPointsModel,
		PointsRecordsModel: pointsRecordsModel,
		Producer:           prod,
		Consumer:           cons,
		MessageRpc:         messageRpc,
	}

	// Start consumer
	if err := cons.Start(); err != nil {
		panic(err)
	}

	return svcCtx
}
