package samples

import (
	"errors"
	"github.com/go-redis/redis/v8"
	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v8"
	"log"
	"sync"
	"time"
)

func LockMain() {
	mutexName := "mutex:holdthedoor"
	rClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	pool := goredis.NewPool(rClient)

	rLock := redsync.New(pool)
	// construct new mutex,
	// mutex has default retries 32 times with delay function,
	// but in this example I wanted to create my own retry
	mutex := rLock.NewMutex(mutexName,
		redsync.WithExpiry(20*time.Second),
		redsync.WithTries(1))
	wg := sync.WaitGroup{}
	wg.Add(1)
	go process(mutex, &wg)
	wg.Add(1)
	go process(mutex, &wg)
	wg.Add(1)
	go process(mutex, &wg)

	wg.Wait()
}

func process(mutex *redsync.Mutex, wg *sync.WaitGroup) {
	defer wg.Done()

	var e error
	for e = errors.New("fake"); e != nil; e = mutex.Lock() {
		log.Print("failed to obtain lock trying to sleep 500ms before retrying...")
		time.Sleep(500 * time.Millisecond)
	}

	// simulate heavy works
	log.Print("obtained the lock and doing something....")
	time.Sleep(3 * time.Second)

	log.Print("Done doing something...")
	mutex.Unlock()
}
