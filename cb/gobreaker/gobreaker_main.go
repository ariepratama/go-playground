package main

import (
	"fmt"
	"github.com/sony/gobreaker"
	"io/ioutil"
	"net/http"
	"time"
)
var cb *gobreaker.CircuitBreaker

func init() {
	var settings gobreaker.Settings
	settings.Name = "TEST"
	settings.MaxRequests = 30
	settings.Interval = time.Minute * time.Duration(1)
	settings.OnStateChange = func(name string, from gobreaker.State, to gobreaker.State) {
		fmt.Printf("%s Changed state from=%s to=%s\n", name, from, to)
	}
	settings.Timeout = time.Minute * time.Duration(3)
	cb = gobreaker.NewCircuitBreaker(settings)
}


func tryGet(url string) {
	body, err := cb.Execute(func() (interface{}, error) {
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		return body, nil
	})

	if err != nil {
		fmt.Printf("err=%s\n", err)
	}

	fmt.Printf("body=%s\n", body)
}

func main() {
	for i := 0; i < 50; i++{
		tryGet("api.github.io")
	}
}