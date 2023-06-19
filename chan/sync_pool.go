package main

import (
	"fmt"
	"sync"
)

/*
首先，sync.Pool 是协程安全的，这对于使用者来说是极其方便的。使用前，设置好对象的 New 函数，
用于在 Pool 里没有缓存的对象时，创建一个。之后，在程序的任何地方、任何时候仅通过 Get()、Put() 方法就可以取、还对象。
*/

var pool *sync.Pool

type Person struct {
	Name string
}

func initPool() {
	pool = &sync.Pool{
		New: func() interface{} {
			fmt.Println("Creating a new Person")
			return new(Person)
		},
	}
}

func main() {
	initPool()

	p := pool.Get().(*Person)
	fmt.Println("首次从 pool 里获取：", p)

	p.Name = "first"
	fmt.Printf("设置 p.Name = %s\n", p.Name)

	pool.Put(p)

	fmt.Println("Pool 里已有一个对象：&{first}，调用 Get: ", pool.Get().(*Person))
	fmt.Println("Pool 没有对象了，调用 Get: ", pool.Get().(*Person))
}
