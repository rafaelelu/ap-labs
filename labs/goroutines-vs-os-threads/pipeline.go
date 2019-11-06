package main

import(
	"time"
	"fmt"
)

const layout = "15:04:05.000000"

func main() {
	nPipes := 1000 // The limit of pipes is around 1000000
	var channels = make([]chan string, 0)
	for i := 0; i < nPipes; i++ {
		channels = append(channels, make(chan string))
	}
	go firstPipe(channels[0])
	for i := 1; i < nPipes; i++ {
		go middlePipe(channels[i-1], channels[i])
	}
	startTime, _ := time.Parse(layout, <-channels[nPipes-1])
	diff := time.Since(startTime)
	fmt.Println("The time it took the message to go through all goroutines and channels was:", diff)
}

func firstPipe(chOut chan string) {
	chOut <- time.Now().Format(layout)
	close(chOut)
}

func middlePipe(chIn chan string, chOut chan string) {
	t := <-chIn
	chOut <- t
	close(chOut)
}