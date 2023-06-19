package main

import (
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	// 推荐用法：Add 和 Wait 不要并发执行
	wg.Add(2)
	go func() {
		wg.Done()
	}()
	go func() {
		wg.Done()
	}()

	wg.Wait()
	fmt.Println("main quit")
}
