package exampleoc2

import (
	"context"
	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/components/cqrs"
	"github.com/ThreeDotsLabs/watermill/message"
)

type OrderCreateEvent struct {
	Id string `json:"id"`
}
type OnOCHandler struct {
	eventBus *cqrs.EventBus
}

func (o OnOCHandler) NewEvent() interface{} {
	//TODO implement me
	panic("implement me")
}

func (o OnOCHandler) Handle(ctx context.Context, cmd interface{}) error {
	//TODO implement me
	panic("implement me")
}

func (o OnOCHandler) HandlerName() string {
	return ""
}

func (o OnOCHandler) NewCommand() interface{} {
	return &OrderCreateEvent{}
}

type OnOrderCreatedCreateReservation struct {
	commandBus *cqrs.CommandBus
}

func (o OnOrderCreatedCreateReservation) HandlerName() string {
	return ""
}

func (OnOrderCreatedCreateReservation) NewEvent() interface{} {
	return nil
}

func init() {
	logger := watermill.NewStdLogger(false, false)
	router, _ := message.NewRouter(message.RouterConfig{}, logger)
	marshaller := cqrs.JSONMarshaler{}
	commandPublisher, _ := kafka.NewPublisher(
		kafka.PublisherConfig{
			Brokers:   []string{"localhost:9092"},
			Marshaler: kafka.DefaultMarshaler{},
		},
		logger)
	commandSubscriber, _ := kafka.NewSubscriber(
		kafka.SubscriberConfig{
			Brokers:               []string{"localhost:9092"},
			Unmarshaler:           kafka.DefaultMarshaler{},
			OverwriteSaramaConfig: kafka.DefaultSaramaSubscriberConfig(),
			ConsumerGroup:         "test_consumer_group",
		},
		logger)

	cqrs.NewFacade(
		cqrs.FacadeConfig{
			CommandHandlers: func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.CommandHandler {
				return []cqrs.CommandHandler{
					OnOCHandler{eventBus: eventBus},
				}
			},
			CommandsSubscriberConstructor: func(handlerName string) (message.Subscriber, error) {
				return commandSubscriber, nil
			},
			CommandsPublisher: commandPublisher,
			EventHandlers: func(commandBus *cqrs.CommandBus, eventBus *cqrs.EventBus) []cqrs.EventHandler {
				return []cqrs.EventHandler{
					OnOCHandler{eventBus: eventBus},
				}
			},
			EventsPublisher:       commandPublisher,
			Router:                router,
			CommandEventMarshaler: marshaller,
			Logger:                logger,
		})
	if err := router.Run(context.Background()); err != nil {
		panic(err)
	}
}
