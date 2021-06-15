package main

import "sync"

func main() {
	var wg sync.WaitGroup
	var a string

	wg.Add(1)

	go func() {
		a = "hello world"
		wg.Done()
	}()

	wg.Wait()
	println(a)
}
