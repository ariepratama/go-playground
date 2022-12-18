package samples

import (
	"context"
	"github.com/go-redis/redis/v8"
)

func PubSubMain() {
	ctx := context.Background()
	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	channelName := "channel"

	rClient.Publish(ctx, channelName, "msg1")

}
