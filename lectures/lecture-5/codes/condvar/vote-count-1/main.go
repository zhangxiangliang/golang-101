package main

import (
	"math/rand"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())

	var count int = 0
	var finished int = 0

	for i := 0; i < 10; i++ {
		go func() {
			vote := requestVote()

			if vote {
				count = count + 1
			}

			finished = finished + 1
		}()
	}

	for count < 5 && finished != 10 {
		// wait
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
