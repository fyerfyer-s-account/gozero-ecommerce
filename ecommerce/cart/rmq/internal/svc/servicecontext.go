package svc

import (
    "context"
    "fmt"
    "log"
    "time"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/consumer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/internal/config"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rmq/producer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/cart/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/broker"
    rmqconfig "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/config"
    "github.com/streadway/amqp"
    "github.com/zeromicro/go-zero/core/stores/sqlx"
)

type ServiceContext struct {
    Config  config.Config
    Broker  broker.Broker
    Channel *amqp.Channel

    // Models
    CartItemsModel    model.CartItemsModel
    CartStatsModel    model.CartStatisticsModel

    // RMQ Components
    Producer *producer.CartProducer
    Consumer *consumer.CartConsumer
}

func NewServiceContext(c config.Config) *ServiceContext {
    log.Println("Initializing ServiceContext...")

    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    // Initialize database connection
    conn := sqlx.NewMysql(c.Mysql.DataSource)

    // Initialize models
    cartItemsModel := model.NewCartItemsModel(conn, c.CacheRedis)
    cartStatsModel := model.NewCartStatisticsModel(conn, c.CacheRedis)

    // Initialize RabbitMQ broker
    rmqBroker := initializeBroker(ctx, &c)
    ch, err := rmqBroker.Channel()
    if err != nil {
        log.Fatalf("Failed to create channel: %v", err)
    }

    // Setup exchanges and queues
    if err := setupRabbitMQ(ch, &c); err != nil {
        log.Fatalf("Failed to setup RabbitMQ: %v", err)
    }

    // Initialize producer and consumer with updated structure
    prod := producer.NewCartProducer(ch, "cart.events")
    cons := consumer.NewCartConsumer(
        ch,
        cartItemsModel,
        cartStatsModel,
    )

    return &ServiceContext{
        Config:          c,
        Broker:         rmqBroker,
        Channel:        ch,
        CartItemsModel: cartItemsModel,
        CartStatsModel: cartStatsModel,
        Producer:       prod,
        Consumer:       cons,
    }
}

func initializeBroker(ctx context.Context, c *config.Config) broker.Broker {
    rmqConfig := convertToEventbusConfig(c)
    rmqBroker := broker.NewAMQPBroker(rmqConfig)

    if err := rmqBroker.Connect(ctx); err != nil {
        log.Fatalf("Failed to connect to RabbitMQ: %v", err)
    }

    return rmqBroker
}

func setupRabbitMQ(ch *amqp.Channel, c *config.Config) error {
    // Setup exchanges
    for _, e := range c.RabbitMQ.Exchanges {
        if err := ch.ExchangeDeclare(
            e.Name,
            e.Type,
            e.Durable,
            e.AutoDelete,
            e.Internal,
            e.NoWait,
            nil,
        ); err != nil {
            return fmt.Errorf("failed to declare exchange %s: %w", e.Name, err)
        }
    }

    // Setup queues and bindings
    for _, q := range c.RabbitMQ.Queues {
        queue, err := ch.QueueDeclare(
            q.Name,
            q.Durable,
            q.AutoDelete,
            q.Exclusive,
            q.NoWait,
            nil,
        )
        if err != nil {
            return fmt.Errorf("failed to declare queue %s: %w", q.Name, err)
        }

        // Setup bindings
        for _, b := range q.Bindings {
            if err := ch.QueueBind(
                queue.Name,
                b.RoutingKey,
                b.Exchange,
                b.NoWait,
                nil,
            ); err != nil {
                return fmt.Errorf("failed to bind queue %s to exchange %s: %w", 
                    queue.Name, b.Exchange, err)
            }
        }
    }

    return nil
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