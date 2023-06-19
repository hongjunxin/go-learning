package main

import (
	"fmt"
	"sync"
	"time"
)

// 导致 panic
// 推荐用法：Add() 和 Wait() 只在同一个 goroutine 中使用
func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		time.Sleep(time.Second)
		wg.Done()
		wg.Add(1)
		fmt.Println("goroutine quit")
	}()
	wg.Wait()
	fmt.Println("main quit")
}
