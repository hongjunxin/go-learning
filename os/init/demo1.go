package main

import "fmt"

func main() {
	a, b, c := 1, "s", 1.01
	fmt.Printf("a=%v, b=%v, c=%v\n", a, b, c)

	// 编译不通过，应该是 [5]int
	//d := [5]{1,2,3}
}
