package main

import "fmt"

type myType struct {
	name string
	age  int
}

func main() {
	ch := make(chan string, 1)
	ch <- "Hi"
	close(ch)

	for {
		elem, ok := <-ch
		fmt.Printf("elem: %v, ok: %v\n", elem, ok)
		if !ok {
			break
		}
	}

	ch2 := make(chan myType, 1)
	ch2 <- myType{name: "Hi", age: 10}
	close(ch2)

	for {
		elem, ok := <-ch2
		fmt.Printf("elem: %v, ok: %v\n", elem, ok)
		if !ok {
			break
		}
	}
}
