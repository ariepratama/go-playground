package example1

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog/log"
	"sync"
	"time"
)

func process(messages <-chan *message.Message) {
	for msg := range messages {
		log.Info().Msgf("consuming message %s", msg.Payload)
		msg.Ack()
	}
}

func publishMessage(wg *sync.WaitGroup, publisher message.Publisher, threadNo int) {
	msgPayload := fmt.Sprintf("hello world! %d: %s", threadNo, time.Now())
	log.Info().Msgf("publishing message with content: %s", msgPayload)
	msg := message.NewMessage(watermill.NewUUID(), []byte(msgPayload))
	if err := publisher.Publish("example.topic", msg); err != nil {
		panic(err)
	}
	wg.Done()
}

func InitSubscriber() {
	subscriberConfig := kafka.DefaultSaramaSubscriberConfig()
	subscriberConfig.Consumer.Offsets.Initial = sarama.OffsetOldest

	subscriber, err := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               []string{"localhost:9092"},
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: subscriberConfig,
			ConsumerGroup:         "test_consumer_group",
		},
		watermill.NewStdLogger(false, false))

	if err != nil {
		panic(err)
	}

	messages, err := subscriber.Subscribe(context.Background(), "example.topic")

	go process(messages)

	if err != nil {
		panic(err)
	}
}

func InitPublisher() {

	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{"localhost:9092"},
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	defer func() {
		log.Debug().Msgf("waiting 5 seconds before exiting the init publisher")
		time.Sleep(5 * time.Second)
		publisher.Close()
	}()
	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("initiating publishing message now")
	numWorkers := 10

	workersGroup := &sync.WaitGroup{}
	workersGroup.Add(10)
	for i := 0; i < numWorkers; i++ {
		go publishMessage(workersGroup, publisher, i)
	}

	workersGroup.Wait()

}
