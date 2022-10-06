package exampleoc

import (
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/rs/zerolog/log"
)

const (
	Retry1MinuteTopic string = "retry.1.minute"
)

func KafkaRetryMiddleware(publisher message.Publisher) func(handlerFunc message.HandlerFunc) message.HandlerFunc {
	return func(handlerFunc message.HandlerFunc) message.HandlerFunc {
		return func(msg *message.Message) ([]*message.Message, error) {
			_, err := handlerFunc(msg)
			// ack the message no matter what happen in the end
			defer msg.Ack()

			// purposely retry if not error
			if err == nil {
				log.Warn().Msg("retrying with kafka...")
				_ = publisher.Publish(Retry1MinuteTopic, msg.Copy())
			}
			return nil, nil
		}
	}
}
