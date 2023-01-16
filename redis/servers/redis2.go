package servers

import (
	"fmt"
	"github.com/go-playground/redis/samples"
	"net/http"
)

var (
	redisRepository2 samples.Repository
)

func redis2(responseWriter http.ResponseWriter, httpRequest *http.Request) {
	key := httpRequest.URL.Query().Get("key")
	res := redisRepository2.Get(key)
	responseWriter.Write([]byte(fmt.Sprintf("[Redis2] Loading from cache results: %s", res)))
}

func InitServerRedis2() {
	redisRepository2 = samples.NewAggCachedRepository()
	http.HandleFunc("/redis2", redis2)
	http.ListenAndServe(":8082", nil)
}
