package samples

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

type (
	Repository interface {
		Get(key string) string
		Set(key, value string)
	}
	dbRepository struct {
		db map[string]string
	}
	cachedRepository struct {
		redisClient *redis.Client
		repository  *dbRepository
	}
)

func NewDbRepository() Repository {
	repo := &dbRepository{db: make(map[string]string)}
	// initialize the db
	for i := 0; i < 1000; i++ {
		key := fmt.Sprintf("%d", i)
		repo.Set(key, key)
	}
	return repo
}

func (r dbRepository) Get(key string) string {
	fmt.Printf("Getting %v from dbrepository...\n", key)
	return r.db[key]
}

func (r dbRepository) Set(key, value string) {
	fmt.Printf("Setting %v from dbrepository...\n", key)
	r.db[key] = value
}

// Get returns the value for the given key from the cache if it exists,
// otherwise return empty string.
func (r cachedRepository) Get(key string) string {
	cmd := r.redisClient.Get(context.Background(), key)
	if cmd.Err() != nil {
		return ""
	}
	fmt.Printf("=======Getting %v from cacherepository=======\n", key)
	return cmd.Val()
}

func (r cachedRepository) Set(key, value string) {
	r.redisClient.Set(context.Background(), key, value, time.Minute)
}

func (r cachedRepository) populateCache(finished chan int) {
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * 50)
		key := fmt.Sprintf("%d", i)

		fmt.Printf("Populating key %s\n", key)
		r.Set(key, key)

	}
	fmt.Println("Finished populating cache")
	finished <- 1
}

func (r cachedRepository) testCacheKey(finished chan int) {
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
	repository := &dbRepository{db: make(map[string]string)}
	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	cachedRepo := &cachedRepository{
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
