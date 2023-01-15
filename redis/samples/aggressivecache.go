package samples

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type (
	// AggCachedRepository is a repository that will not do lazy loading
	AggCachedRepository struct {
		redisClient *redis.Client
		repository  *DbRepository
	}
)

// Get returns the value for the given key from the cache if it exists,
// otherwise return empty string.
func (r AggCachedRepository) Get(key string) string {
	cmd := r.redisClient.Get(context.Background(), key)
	if cmd.Err() != nil {
		return ""
	}
	fmt.Printf("=======Getting %v from cacherepository=======\n", key)
	return cmd.Val()
}

func (r AggCachedRepository) Set(key, value string) {
	r.redisClient.Set(context.Background(), key, value, time.Minute)
}

func (r AggCachedRepository) populateCache(finished chan int) {
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 50)
		key := fmt.Sprintf("%d", i)

		fmt.Printf("Populating key %s\n", key)
		r.Set(key, key)

	}
	fmt.Println("Finished populating cache")
	finished <- 1
}

func (r AggCachedRepository) testCacheKey(finished chan int) {
	for i := 100; i > 0; i-- {
		time.Sleep(time.Millisecond * 75)
		key := fmt.Sprintf("%d", i)
		v := r.Get(key)
		if v == "" {
			fmt.Printf("key %s is not ready yet\n", key)
		}
	}
	fmt.Println("Finished testing cache")
	finished <- 2
}

func AggressiveCacheMain() {
	repository := &DbRepository{db: make(map[string]string)}
	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	cachedRepo := &AggCachedRepository{
		rClient,
		repository}
	val := cachedRepo.Get("1")
	if val != "" {
		panic("not expected")
	}

	populatingFinished := make(chan int, 2)
	go cachedRepo.populateCache(populatingFinished)
	go cachedRepo.testCacheKey(populatingFinished)
	<-populatingFinished
	print("clearing up the keys...")
	keys := rClient.Keys(context.Background(), "*").Val()
	rClient.Del(context.Background(), keys...)
}
