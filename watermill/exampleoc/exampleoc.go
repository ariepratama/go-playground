package exampleoc

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/ThreeDotsLabs/watermill/message/router/middleware"
	"github.com/ThreeDotsLabs/watermill/message/router/plugin"
	zerolog "github.com/rs/zerolog/log"
	"log"
	"sync"
	"time"
)

var (
	logger = watermill.NewStdLogger(false, false)
)

type structHandler struct {
	// we can add some dependencies here
}

func (s structHandler) Handler(msg *message.Message) ([]*message.Message, error) {
	log.Println("structHandler received message", msg.UUID, string(msg.Payload))

	msg = message.NewMessage(watermill.NewUUID(), []byte("message produced by structHandler"))
	return message.Messages{msg}, nil
}

func publishMessage(wg *sync.WaitGroup, publisher message.Publisher, topic string) {
	msgPayload := fmt.Sprintf("order create %s", time.Now())
	msg := message.NewMessage(watermill.NewUUID(), []byte(msgPayload))
	if err := publisher.Publish(topic, msg); err != nil {
		panic(err)
	}
	zerolog.Debug().Msgf("done publishing message to topic %s", topic)
	wg.Done()
}

func InitPublisher() {

}

func InitOrderCreateConsumer() {
	router, _ := message.NewRouter(message.RouterConfig{}, logger)

	subscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	subscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest
	kafkaSubscriber, _ := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               []string{"localhost:9092"},
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: subscriberConfig,
			ConsumerGroup:         "test_consumer_group",
		},
		logger)
	kafkaPublisher, _ := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{"localhost:9092"},
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger)
	// kill router when SIGTERM is issued
	router.AddPlugin(plugin.SignalsHandler)
	router.AddMiddleware(
		middleware.CorrelationID,
		//middleware.Retry{
		//	MaxRetries:      3,
		//	InitialInterval: time.Second,
		//	Logger:          logger,
		//}.Middleware,
		KafkaRetryMiddleware(kafkaPublisher),
		middleware.Recoverer,
	)
	router.AddHandler(
		"order_create_handler",
		"order_create_topic",
		kafkaSubscriber,
		"order_create_result_topic",
		kafkaPublisher,
		structHandler{}.Handler,
	)
	zerolog.Debug().Msg("starting...")
	// run in different goroutine
	go func() {
		if err := router.Run(context.Background()); err != nil {
			panic(err)
		}
	}()

	numMessages := 10
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(numMessages)
	for i := 0; i < numMessages; i++ {
		go publishMessage(waitGroup, kafkaPublisher, "order_create_topic")
	}

	waitGroup.Wait()

}
