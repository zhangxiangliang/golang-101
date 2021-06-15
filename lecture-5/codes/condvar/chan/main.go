package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var c chan bool = make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			c <- requestVote()
		}()
	}

	var count int = 0
	var finished int = 0

	for {
		var vote bool = <-c
		if vote {
			count = count + 1
		}
		finished = finished + 1

		if count > 5 {
			break
		}

		if finished == 10 {
			break
		}
	}

	if count >= 5 {
		println("received 5+ votes!")
	} else {
		println("lost")
	}
}

func requestVote() bool {
	time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
	return rand.Int()%2 == 0
}
