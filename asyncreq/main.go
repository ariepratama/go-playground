// simple example on how to use asyncreq

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ariepratama/asyncreq/core"
	"github.com/go-redis/redis/v8"
	"io"
	"io/ioutil"
	"net/http"
)

var (
	redisClient *redis.Client
	postHandler *core.RedisPostHandler
	getHandler  *core.RedisGetHandler
)

func callback(request core.PostRequest) {
	fmt.Sprintf("callback called %s...\n", request.Payload)
}

func processReq(ctx context.Context, request core.PostRequest, callback core.ProcessReqFunCallback) {
	fmt.Sprintf("processing request %s...\n", request.Payload)
}

func init() {
	redisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	postHandler = &core.RedisPostHandler{
		RedisClient: redisClient,
		OnPostRequest: core.OnPostRequest{
			Callback:      callback,
			ProcessReqFun: processReq,
		},
	}
	getHandler = &core.RedisGetHandler{
		RedisClient: redisClient,
	}
}

func getRoot(responseWritter http.ResponseWriter, request *http.Request) {
	io.WriteString(responseWritter, "hello there!")
}

func asyncRequestRouter(responseWriter http.ResponseWriter, request *http.Request) {
	if http.MethodPost == request.Method {
		postAsyncRequest(responseWriter, request)
		return
	}

	if http.MethodGet == request.Method {
		getAsyncRequest(responseWriter, request)
		return
	}

	responseWriter.WriteHeader(http.StatusBadRequest)
	io.WriteString(responseWriter, "http method not allowed!\n")
}

func getAsyncRequest(responseWriter http.ResponseWriter, request *http.Request) {
	defer responseWriter.Header().Set("Content-Type", "application/json")

	requestId := request.URL.Query().Get("request_id")

	if requestId == "" {
		responseWriter.WriteHeader(http.StatusBadRequest)
		io.WriteString(responseWriter, "{\"error\": \"request id should not be empty\"}")
		return
	}
	getResponse := getHandler.Do(requestId)
	responseWriter.WriteHeader(http.StatusOK)
	json.NewEncoder(responseWriter).Encode(getResponse)
}

func postAsyncRequest(responseWriter http.ResponseWriter, request *http.Request) {
	defer responseWriter.Header().Set("Content-Type", "application/json")

	contentBytes, err := ioutil.ReadAll(request.Body)
	if err != nil {
		responseWriter.WriteHeader(http.StatusInternalServerError)
		io.WriteString(responseWriter, "{\"error\": \"unexpected error happened\"}")
		return
	}

	postResponse := postHandler.Do(core.PostRequest{
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
