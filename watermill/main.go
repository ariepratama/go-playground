package main

import (
	"fmt"
	"time"
	"watermill/exampleoc"
)

func main() {
	fmt.Print("start")
	exampleoc.InitOrderCreateConsumer()
	exampleoc.InitPublisher()
	time.Sleep(10 * time.Second)
	fmt.Print("finish")
}
