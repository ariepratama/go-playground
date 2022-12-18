package samples

import (
	"context"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"log"
	"time"
)

const (
	taskQueueName  = "taskqueue:sample1"
	maxWorkerRetry = 3
)

// TaskQueueMain simulates one fast publisher and other parallel workers trying to keep up
func TaskQueueMain() {
	nWorkers := 2
	publishFinished := make(chan bool)
	taskQueueInitialized := make(chan bool)
	workersFinished := make(chan bool, nWorkers)
	// start an asyncPublisher
	go startPublisher(100, nWorkers, publishFinished, taskQueueInitialized)
	// wait for signal from publisher to start consumptions
	<-taskQueueInitialized
	// asynchronously start parallel workers
	for i := 0; i < nWorkers; i++ {
		go startWorker(workersFinished)
	}
	// wait until all messages have been published
	<-publishFinished
	// wait until all workers have given up
	<-workersFinished
}

func defaultRedisClient() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
}

// startPublisher that will generate 4 tasks per 1s
func startPublisher(nTasks int, nWorkers int, publishFinished chan bool, taskQueueInitialized chan bool) {
	rClient := defaultRedisClient()
	initialized := false
	generatedTasks := 0
	for {
		taskUuid, _ := uuid.NewUUID()
		taskId := taskUuid.String()
		log.Printf("Inserting random task id %s", taskId)

		rClient.LPush(context.Background(), taskQueueName, taskId)

		if generatedTasks >= nWorkers && !initialized {
			taskQueueInitialized <- true
			initialized = !initialized
			log.Print("Sending signal to start workers...")

		}

		generatedTasks++
		time.Sleep(25 * time.Millisecond)
		if generatedTasks > nTasks {
			publishFinished <- true
		}
	}
}

func startWorker(consumeFinished chan bool) {
	tag := "CONSUMER"
	rClient := defaultRedisClient()
	retryCount := 0
	for {
		popResult := rClient.RPop(context.Background(), taskQueueName)
		if popResult.Err() != nil {
			log.Printf("stopping consumer because of error %v", popResult.Err())
			if retryCount > maxWorkerRetry {
				log.Printf("There is nothing else in the queue, giving up")
				consumeFinished <- true
				break
			}
			retryCount++
			time.Sleep(500 * time.Millisecond)
			continue
		}

		countResult := rClient.LLen(context.Background(), taskQueueName)
		if countResult.Err() != nil {
			log.Printf("stopping consumer because count failed error = %v", countResult.Err())
			consumeFinished <- true
			break
		}
		if countResult.Val() <= 0 {
			log.Printf("There is nothing else in the queue, retry=%d", retryCount)
			if retryCount > maxWorkerRetry {
				log.Printf("There is nothing else in the queue, giving up")
				consumeFinished <- true
				break
			}
			retryCount++
			time.Sleep(500 * time.Millisecond)
			continue
		}

		taskId := popResult.Val()
		log.Printf("[%s] doing task %s", tag, taskId)
		time.Sleep(500 * time.Millisecond)
		log.Printf("[%s]finished doing task %v", tag, taskId)
	}
}
