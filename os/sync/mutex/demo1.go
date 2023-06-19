package main

import (
	"fmt"
	"sync"
)

var ch = make(chan struct{})

func main() {
	var mu sync.Mutex
	mu.Lock()
	fmt.Println("main mu.Lock() ok")
	go test1(mu) // 作为值传递，会将锁的状态一起复制，比如这里已经被 lock
	mu.Unlock()
	fmt.Println("main mu.Unlock() ok")
	<-ch
}

func test1(m sync.Mutex) {
	m.Unlock()
	fmt.Println("test1 m.Unlock ok")
	m.Lock()
	fmt.Println("test1 m.Lock() ok")
	m.Unlock()
	fmt.Println("test1 m.Unlock() ok")
	ch <- struct{}{}
}
