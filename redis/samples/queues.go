package samples

import (
	"context"
	"github.com/go-redis/redis/v9"
	"log"
	"math/rand"
	"sync"
	"time"
)

func QueueMain() {
	ctx := context.Background()

	key := "tasks:queue"

	populatingDone := make(chan bool)
	// add new element at the head of list asynchronously
	go func() {
		rClient := redis.NewClient(&redis.Options{
			Addr:     "localhost:6379",
			Password: "",
			DB:       0,
		})
		rClient.LPush(ctx, key, 300)
		rClient.LPush(ctx, key, 301)
		rClient.LPush(ctx, key, 302)
		populatingDone <- true
	}()

	// wait until populating is done, then continue to pop
	<-populatingDone
	// consuming element from the bottom of the queue asynchronously
	log.Print("Popping from list....")

	// will always results in
	// 	2022/12/11 22:47:56 consuming task 302 after sleeping for x
	//	2022/12/11 22:47:58 consuming task 301 after sleeping for x
	//	2022/12/11 22:47:58 consuming task 300 after sleeping for x
	wg := sync.WaitGroup{}
	// wait group add should be outside of goroutine
	// such that the counter will be correct
	// if there's wait group add inside goroutine, the wait might not detect it
	// because on Wait() called, the counter has not been added yet
	wg.Add(1)
	go doWork(key, &wg)
	wg.Add(1)
	go doWork(key, &wg)
	wg.Add(1)
	go doWork(key, &wg)

	wg.Wait()
	log.Print("Done consuming all tasks, terminating")
}

func doWork(key string, wg *sync.WaitGroup) {
	defer wg.Done()
	ctx := context.Background()
	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	randWaitSec := rand.Intn(3) + 1
	log.Printf("Sleeping for %d second...", randWaitSec)
	time.Sleep(time.Duration(randWaitSec) * time.Second)
	result := rClient.LPop(ctx, key)
	log.Printf("consuming task %s after sleeping for %d", result.Val(), randWaitSec)
}
