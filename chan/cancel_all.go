package main

import (
	"fmt"
	"sync"
	"time"
)

func worker(wg *sync.WaitGroup, cannel chan bool, tag int) {
	defer wg.Done()

	for {
		select {
		default:
			fmt.Printf("hello %v\n", tag)
		case <-cannel:
			// do some clean
			return
		}
	}
}

func main() {
	cancel := make(chan bool)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go worker(&wg, cancel, i)
	}

	time.Sleep(time.Second)
	close(cancel) // let all worker quit
	wg.Wait()
}
