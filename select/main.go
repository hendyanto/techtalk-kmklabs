package main

import (
	"fmt"
	"time"
)

func main() {
	channel_1 := make(chan string)
	channel_2 := make(chan string)

	go func() {
		time.Sleep(time.Second * 2)
		channel_1 <- "one"
	}()

	go func() {
		time.Sleep(time.Second * 1)
		channel_2 <- "two"
	}()

	for i := 0; i < 2; i++ {
		select {
		case message1 := <-channel_1:
			fmt.Println("message received", message1)
		case message2 := <-channel_2:
			fmt.Println("message received", message2)
		}
	}
}
