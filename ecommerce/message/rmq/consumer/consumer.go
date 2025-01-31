package consumer

import (
    "context"

    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rpc/model"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/message/rmq/consumer/handlers"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/consumer"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/eventbus/middleware"
    "github.com/fyerfyer/gozero-ecommerce/ecommerce/pkg/zerolog"
    "github.com/streadway/amqp"
)

type MessageConsumer struct {
    logger                *zerolog.Logger
    channel              *amqp.Channel
    eventMessageHandler   *handlers.EventMessageHandler
    templateHandler      *handlers.TemplateMessageHandler
    batchHandler         *handlers.BatchMessageHandler
    paymentSuccessHandler *handlers.PaymentSuccessHandler
    paymentFailedHandler  *handlers.PaymentFailedHandler
    consumers            []*consumer.Consumer
}

func NewMessageConsumer(
    ch *amqp.Channel,
    messagesModel model.MessagesModel,
    messageSendsModel model.MessageSendsModel,
    templatesModel model.MessageTemplatesModel,
    settingsModel model.NotificationSettingsModel,
) *MessageConsumer {
    return &MessageConsumer{
        logger:  zerolog.GetLogger(),
        channel: ch,
        eventMessageHandler: handlers.NewEventMessageHandler(
            messagesModel,
            messageSendsModel,
            templatesModel,
            settingsModel,
        ),
        templateHandler: handlers.NewTemplateMessageHandler(
            messagesModel,
            messageSendsModel,
            templatesModel,
            settingsModel,
        ),
        batchHandler: handlers.NewBatchMessageHandler(
            messagesModel,
            messageSendsModel,
            templatesModel,
            settingsModel,
        ),
        paymentSuccessHandler: handlers.NewPaymentSuccessHandler(
            messagesModel,
            messageSendsModel,
            templatesModel,
            settingsModel,
        ),
        paymentFailedHandler: handlers.NewPaymentFailedHandler(
            messagesModel,
            messageSendsModel,
            templatesModel,
            settingsModel,
        ),
    }
}

func (c *MessageConsumer) Start(ctx context.Context) error {
    // Create event message consumer
    eventConsumer, err := consumer.NewConsumer(
        c.channel,
        c.eventMessageHandler.Handle,
        consumer.WithQueue("message.event"),
        consumer.WithExchange("message.events"),
        consumer.WithRoutingKey("message.event.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("message-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create template message consumer
    templateConsumer, err := consumer.NewConsumer(
        c.channel,
        c.templateHandler.Handle,
        consumer.WithQueue("message.template"),
        consumer.WithExchange("message.events"),
        consumer.WithRoutingKey("message.template.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("message-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create batch message consumer
    batchConsumer, err := consumer.NewConsumer(
        c.channel,
        c.batchHandler.Handle,
        consumer.WithQueue("message.batch"),
        consumer.WithExchange("message.events"),
        consumer.WithRoutingKey("message.batch.*"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("message-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create payment success consumer
    paymentSuccessConsumer, err := consumer.NewConsumer(
        c.channel,
        c.paymentSuccessHandler.Handle,
        consumer.WithQueue("message.payment.success"),
        consumer.WithExchange("payment.events"),
        consumer.WithRoutingKey("payment.success"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("message-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Create payment failed consumer
    paymentFailedConsumer, err := consumer.NewConsumer(
        c.channel,
        c.paymentFailedHandler.Handle,
        consumer.WithQueue("message.payment.failed"),
        consumer.WithExchange("payment.events"),
        consumer.WithRoutingKey("payment.failed"),
        consumer.WithAutoAck(false),
        consumer.WithMiddlewares(
            middleware.Recovery,
            middleware.Logging,
            middleware.Retry(nil),
            middleware.Tracing("message-consumer"),
        ),
    )
    if err != nil {
        return err
    }

    // Start all consumers
    for _, c := range []struct {
        name     string
        consumer *consumer.Consumer
    }{
        {"event", eventConsumer},
        {"template", templateConsumer},
        {"batch", batchConsumer},
        {"payment.success", paymentSuccessConsumer},
        {"payment.failed", paymentFailedConsumer},
    } {
        if err := c.consumer.Start(ctx); err != nil {
            return err
        }
    }

    c.consumers = []*consumer.Consumer{
        eventConsumer,
        templateConsumer,
        batchConsumer,
        paymentSuccessConsumer,
        paymentFailedConsumer,
    }

    return nil
}

func (c *MessageConsumer) Stop() error {
    for _, consumer := range c.consumers {
        if err := consumer.Stop(); err != nil {
            c.logger.Error(context.Background(), "Failed to stop consumer", err, nil)
        }
    }
    return nil
}