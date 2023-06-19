package main

import (
	"context"
	"fmt"
	"sync"
	"time"
)

func workerV2(ctx context.Context, wg *sync.WaitGroup, tag int) error {
	defer wg.Done()

	for {
		time.Sleep(time.Second)
		select {
		default:
			fmt.Printf("worker %v\n", tag)
		case <-ctx.Done():
			fmt.Printf("worker %v, err='%v'\n", tag, ctx.Err().Error())
			return ctx.Err()
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)

	var wg sync.WaitGroup
	for i := 0; i < 2; i++ {
		wg.Add(1)
		go workerV2(ctx, &wg, i)
	}

	time.Sleep(time.Second * 3)

	// worker quit with err='context canceled'
	cancel()

	// worker quit with err='context deadline exceeded'
	//defer cancel()

	wg.Wait()
}
