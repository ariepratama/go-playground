package samples

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"log"
	"time"
)

func HashMain() {
	ctx := context.Background()
	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	key := "Some Key"
	val := "Some Val"
	status := rClient.Set(ctx, key, val, 59*time.Second)
	if status.Err() != nil {
		log.Printf("Error in setting key=%s, err=%s", key, status.Err())
		panic(status.Err())
	}

	log.Print(fmt.Printf("Setting redis key %s with value %s", key, val))

	getVal := rClient.Get(ctx, key)
	if getVal.Err() != nil {
		panic(getVal.Err())
	}

	log.Print(fmt.Printf("Getting redis key %s, result value=%s", key, getVal.Val()))
}
