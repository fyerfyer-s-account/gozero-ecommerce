package svc

import (
    "context"
    "log"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/consumer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rmq/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/inventory/rpc/model"
    rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/broker"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
    Config config.Config
    Broker broker.Broker

    // Models
    StocksModel       model.StocksModel
    StockLocksModel   model.StockLocksModel
    StockRecordsModel model.StockRecordsModel
    WarehousesModel   model.WarehousesModel

    // RMQ Components
    Producer *producer.InventoryProducer
    Consumer *consumer.InventoryConsumer
}

func NewServiceContext(c config.Config) *ServiceContext {
    log.Println("Initializing ServiceContext...")

    // Connect to MySQL
    log.Println("Connecting to MySQL...")
    conn := sqlx.NewMysql(c.Mysql.DataSource)
    log.Println("Connected to MySQL successfully.")

    // Initialize RabbitMQ broker
    log.Println("Initializing RabbitMQ broker...")
    rmqConfig := convertToEventbusConfig(&c)
    rmqBroker := broker.NewAMQPBroker(rmqConfig)
    log.Println("Connecting to RabbitMQ broker...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    
    if err := rmqBroker.Connect(ctx); err != nil {
        log.Fatalf("Failed to connect to RabbitMQ broker: %v", err)
    }
    log.Println("RabbitMQ broker connected successfully.")

    // Get RMQ channel
    log.Println("Getting RabbitMQ channel...")
    ch, err := rmqBroker.Channel()
    if err != nil {
        log.Fatalf("Failed to get RabbitMQ channel: %v", err)
    }
    log.Println("RabbitMQ channel acquired successfully.")

    // Initialize models
    log.Println("Initializing models...")
    stocksModel := model.NewStocksModel(conn, c.CacheRedis)
    log.Println("StocksModel initialized.")
    stockLocksModel := model.NewStockLocksModel(conn, c.CacheRedis)
    log.Println("StockLocksModel initialized.")
    stockRecordsModel := model.NewStockRecordsModel(conn, c.CacheRedis)
    log.Println("StockRecordsModel initialized.")
    warehousesModel := model.NewWarehousesModel(conn, c.CacheRedis)
    log.Println("WarehousesModel initialized.")

    // Initialize producer and consumer
    log.Println("Initializing producer...")
    prod := producer.NewInventoryProducer(ch, "inventory.events")
    log.Println("InventoryProducer initialized successfully.")
    log.Println("Initializing consumer...")
    cons := consumer.NewInventoryConsumer(
        ch,
        stocksModel,
        stockLocksModel,
        stockRecordsModel,
    )
    log.Println("InventoryConsumer initialized successfully.")

    return &ServiceContext{
        Config: c,
        Broker: rmqBroker,

        StocksModel:       stocksModel,
        StockLocksModel:   stockLocksModel,
        StockRecordsModel: stockRecordsModel,
        WarehousesModel:   warehousesModel,

        Producer: prod,
        Consumer: cons,
    }
}

func convertToEventbusConfig(c *config.Config) *rmqconfig.RabbitMQConfig {
    log.Println("Converting RabbitMQ configuration...")

    exchanges := make([]rmqconfig.ExchangeConfig, len(c.RabbitMQ.Exchanges))
    for i, e := range c.RabbitMQ.Exchanges {
        log.Printf("Configuring exchange: %s", e.Name)
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
        log.Printf("Configuring queue: %s", q.Name)
        bindings := make([]rmqconfig.BindingConfig, len(q.Bindings))
        for j, b := range q.Bindings {
            log.Printf("Configuring binding: Exchange=%s, RoutingKey=%s", b.Exchange, b.RoutingKey)
            bindings[j] = rmqconfig.BindingConfig{
                Exchange:   b.Exchange,
                RoutingKey: b.RoutingKey,
                NoWait:    b.NoWait,
            }
        }
        queues[i] = rmqconfig.QueueConfig{
            Name:       q.Name,
            Durable:    q.Durable,
            AutoDelete: q.AutoDelete,
            Exclusive:  q.Exclusive,
            NoWait:    q.NoWait,
            Bindings:   bindings,
        }
    }

    log.Println("RabbitMQ configuration conversion complete.")
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