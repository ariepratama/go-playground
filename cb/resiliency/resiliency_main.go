package main

import (
	"fmt"
	"github.com/eapache/go-resiliency/breaker"
	"io/ioutil"
	"net/http"
	"time"
)

var cbBreaker *breaker.Breaker

func init() {
	cbBreaker = breaker.New(10, 1, time.Minute*time.Duration(3))
}

func getSome(url string, callback func(interface{})) {
	var b interface{}
	errResult := cbBreaker.Run(func() error {
		resp, err := http.Get(url)
		if err != nil {
			return err
		}

		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		b = body

		return nil
	})

	switch errResult {
	case nil:
		callback(b)
		break
	case breaker.ErrBreakerOpen:
		fmt.Printf("%s\n", errResult)
		break
	default:
		fmt.Printf("What is going on here ? %s\n", errResult)
		break
	}
}

func main() {
	for i := 0; i < 50; i++ {
		getSome("api.github.io", func(body interface{}) {
			fmt.Printf("Call results =%s\n", body)
		})
	}
}
