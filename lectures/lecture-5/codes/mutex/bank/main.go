package main

import (
	"fmt"
	"sync"
	"time"
)

func main() {
	var mu sync.Mutex

	alice := 1000
	bob := 1000
	total := alice + bob

	go func() {
		mu.Lock()
		alice = alice - 1
		bob = bob + 1
		mu.Unlock()
	}()

	go func() {
		mu.Lock()
		alice = alice + 1
		bob = bob - 1
		mu.Unlock()
	}()

	start := time.Now()
	for time.Since(start) < 1*time.Second {
		mu.Lock()
		if alice+bob != total {
			fmt.Printf("observed violation, alice = %v, bob = %v, sum = %v\n", alice, bob, alice+bob)
		}
		mu.Unlock()
	}
}
