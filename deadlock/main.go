package main

import "fmt"

func main() {
	fmt.Println("Hello world!")
	channel_start := make(chan string)
	channel_end := make(chan string)

	channel_start <- "Message"

	<-channel_end
}
