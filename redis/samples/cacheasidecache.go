package samples

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

// https://docs.aws.amazon.com/whitepapers/latest/database-caching-strategies-using-redis/caching-patterns.html
// cache aside is lazy loading, means whenever there's a cache missed, then it will be loaded from the database

type (
	CacheAsideRepository struct {
		redisClient *redis.Client
		repository  Repository
	}
)

func NewCacheAsideRepository() Repository {
	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &CacheAsideRepository{redisClient: rClient, repository: NewDbRepository()}
}

func (repo CacheAsideRepository) Get(key string) string {
	cmd := repo.redisClient.Get(context.Background(), key)
	// the key is not existed
	if redis.Nil == cmd.Err() {
		fmt.Printf("=======Getting %v from dbrepository=======\n", key)
		// simulate getting from database
		time.Sleep(DbLoadTime)
		// get the value from the database
		value := repo.repository.Get(key)
		// set the value to the cache
		repo.Set(key, value)
		return value
	}

	if cmd.Err() != nil {
		return ""
	}

	fmt.Printf("=======Getting %v from cacherepository=======\n", key)
	return cmd.Val()
}

func (repo CacheAsideRepository) Set(key, value string) {
	repo.redisClient.Set(context.Background(), key, value, time.Minute)
}
