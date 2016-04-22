package main

import (
	"fmt"
	"time"
)

func susunKotak(in <-chan string, out chan<- string) {
	for c := range in {
		time.Sleep(time.Second * 1)
		var processed string = c + "[kotak]"
		fmt.Println("\n", processed)
		out <- processed
	}
}

func masukkanBubur(in chan string, out chan string) {
	for c := range in {
		time.Sleep(time.Second * 2)
		var processed string = c + "[bubur]"
		fmt.Println("\n", processed)
		out <- processed
	}
}

func masukkanTopping(in chan string, out chan string) {
	for c := range in {
		time.Sleep(time.Second * 3)
		var processed string = c + "[topping]"
		fmt.Println("\n", processed)
		out <- processed
	}
}

func bungkus(in chan string, out chan string) {
	for c := range in {
		time.Sleep(time.Second * 1)
		var processed string = c + "[bungkus]"
		fmt.Println("\n", processed)
		out <- processed
	}
}

func main() {
	channel_1 := make(chan string, 10)
	channel_2 := make(chan string, 10)
	channel_3 := make(chan string, 10)
	channel_4 := make(chan string, 10)
	channel_finished := make(chan string, 10)

	target := []string{"A", "B", "C", "D", "E", "F", "G", "H", "I", "J"}

	go susunKotak(channel_1, channel_2)

	go masukkanBubur(channel_2, channel_3)

	go masukkanTopping(channel_3, channel_4)

	go bungkus(channel_4, channel_finished)

	for _, t := range target {
		fmt.Println(" ", t, "[pesan]")
		channel_1 <- t + "[pesan]"
	}

	for range target {
		<-channel_finished
	}
}
