// simple example on how to use asyncreq

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ariepratama/asyncreq/asyncreq"
	"github.com/go-redis/redis/v8"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type (
	AsyncService struct {
		postHandler asyncreq.PostHandler
		getHandler  asyncreq.GetHandler
	}
)

var (
	service *AsyncService
)

func onPostRequestCompleted(ctx context.Context, request *asyncreq.PostRequest, response asyncreq.PostResponse) asyncreq.PostResponse {
	log.Println(fmt.Sprintf("callback called %s...\n", request.Payload))
	return response
}

func onPostError(ctx context.Context, err error) asyncreq.PostResponse {
	log.Println("post error...")

	return asyncreq.PostResponse{}
}

func onPostRequest(ctx context.Context, request *asyncreq.PostRequest) asyncreq.PostResponse {
	log.Println(fmt.Sprintf("processing request"))
	time.Sleep(time.Second * 5)
	log.Println(fmt.Sprintf("finished request"))
	return asyncreq.PostResponse{
		IsError:      false,
		ErrorMessage: "",
		RequestId:    "assafds",
	}
}

func init() {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	service = &AsyncService{
		postHandler: asyncreq.NewRedisPostHandler(
			redisClient,
			asyncreq.PostRequestRedisOptions{
				Ttl: time.Second * 30,
			},
			onPostRequest,
			onPostRequestCompleted,
			onPostError,
		),
		getHandler: asyncreq.RedisGetHandler{
			RedisClient: redisClient,
		},
	}
}

func getRoot(responseWritter http.ResponseWriter, request *http.Request) {
	io.WriteString(responseWritter, "hello there!")
}

func asyncRequestRouter(responseWriter http.ResponseWriter, request *http.Request) {
	if http.MethodPost == request.Method {
		postRequestHandler(responseWriter, request)
		return
	}

	if http.MethodGet == request.Method {
		getRequestHandler(responseWriter, request)
		return
	}

	responseWriter.WriteHeader(http.StatusBadRequest)
	io.WriteString(responseWriter, "http method not allowed!\n")
}

func getRequestHandler(responseWriter http.ResponseWriter, request *http.Request) {
	defer responseWriter.Header().Set("Content-Type", "application/json")

	requestId := request.URL.Query().Get("request_id")

	if requestId == "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		io.WriteString(responseWriter, "{\"error\": \"request id should not be empty\"}")
		return
	}
	getResponse := service.getHandler.Do(requestId)
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(getResponse)
}

func postRequestHandler(responseWriter http.ResponseWriter, request *http.Request) {
	defer responseWriter.Header().Set("Content-Type", "application/json")

	contentBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		io.WriteString(responseWriter, "{\"error\": \"unexpected error happened\"}")
		return
	}

	postResponse := service.postHandler.Do(asyncreq.PostRequest{
		Payload: string(contentBytes),
	})

	httpStatus := http.StatusCreated

	if postResponse.IsError {
		httpStatus = http.StatusInternalServerError
	}

	responseWriter.WriteHeader(httpStatus)
	responseWriter.Header().Set("Content-Type", "application/json")
	json.NewEncoder(responseWriter).Encode(postResponse)
}

func main() {
	http.HandleFunc("/", getRoot)
	http.HandleFunc("/request", asyncRequestRouter)

	err := http.ListenAndServe(":9876", nil)
	if err != nil {
		fmt.Sprintf("closing server %s\n", err)
	}
}
