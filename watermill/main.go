package main

import (
	"context"
	"fmt"
	"github.com/Shopify/sarama"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog/log"
	"time"
)

func process(messages <-chan *message.Message) {
	for msg := range messages {
		log.Info().Msgf("consuming message %s", msg.Payload)
		msg.Ack()
	}
}

func publishMessages(publisher message.Publisher) {
	for {
		msg := message.NewMessage(watermill.NewUUID(), []byte(fmt.Sprintf("hello world! %d", time.Now().Second())))
		if err := publisher.Publish("example.topic", msg); err != nil {
			panic(err)
		}
		time.Sleep(time.Second)
	}
}

func initSubscriber() {
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

func initPublisher() {
	defer func() {
		log.Debug().Msgf("waiting 5 seconds before exiting the init publisher")
		time.Sleep(5 * time.Second)
	}()
	publisher, err := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{"localhost:9092"},
			Marshaler: kafka.DefaultMarshaler{},
		},
		watermill.NewStdLogger(false, false),
	)
	if err != nil {
		panic(err)
	}
	log.Debug().Msgf("initiating publishing message now")
	go publishMessages(publisher)

}

func main() {
	fmt.Print("start")

	initSubscriber()
	initPublisher()
	log.Debug().Msgf("Waiting 10s before terminating the program...")
	time.Sleep(10 * time.Second)
	fmt.Print("finish")

}
