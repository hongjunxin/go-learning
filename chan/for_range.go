package main

import (
	"fmt"
)

func main() {
	ch := make(chan int)
	go func() {
		for {
			select {
			case ch <- 0:
			case ch <- 1:
			}
		}
	}()

	// for range 会一直读取 ch，并不会退出，直到 ch 被关闭
	for v := range ch {
		fmt.Println(v)
	}
	fmt.Println("main quit")
}
