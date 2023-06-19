package main

import (
	"container/ring"
	"fmt"
)

func main() {
	r := ring.New(5)
	for i := 0; i < r.Len(); i++ {
		r.Value = i
		r = r.Next()
	}
	r.Do(func(i interface{}) {
		fmt.Printf("%d,", i.(int))
	})
	fmt.Println()
}
