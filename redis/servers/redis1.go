package servers

import (
	"fmt"
	"github.com/go-playground/redis/samples"
	"net/http"
)

var (
	redisRepository samples.Repository
)

func redis1(w http.ResponseWriter, r *http.Request) {
	key := r.URL.Query().Get("key")
	res := redisRepository.Get(key)
	w.Write([]byte(fmt.Sprintf("Loading from cache results: %s", res)))
}

func InitServerRedis1() {
	redisRepository = samples.NewCacheAsideRepository()
	http.HandleFunc("/redis1", redis1)
	http.ListenAndServe(":8081", nil)
}
