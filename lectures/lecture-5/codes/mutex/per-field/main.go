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

	// 重复来回转钱
	go func() {
		for i := 0; i < 1000; i++ {
			mu.Lock()
			alice = alice - 1
			mu.Unlock()

			mu.Lock()
			bob = bob + 1
			mu.Unlock()
		}
	}()

	// 重复来回转钱
	go func() {
		for i := 0; i < 1000; i++ {
			mu.Lock()
			bob = bob - 1
			mu.Unlock()

			mu.Lock()
			alice = alice + 1
			mu.Unlock()
		}
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
