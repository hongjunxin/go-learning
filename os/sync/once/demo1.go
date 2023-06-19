package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	var counter uint32
	var wg sync.WaitGroup
	var once sync.Once
	wg.Add(2)

	go func() {
		defer wg.Done()
		once.Do(func() {
			atomic.AddUint32(&counter, 1)
		})
	}()

	go func() {
		defer wg.Done()
		once.Do(func() {
			atomic.AddUint32(&counter, 2)
		})
	}()

	wg.Wait()
	fmt.Printf("The counter: %d\n", counter)
}
