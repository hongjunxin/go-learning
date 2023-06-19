package main

import "fmt"

func main() {
	input := []int{1, 2, 3, 4, 5}
	fmt.Printf("sum1: %v\n", sum(input...))
	fmt.Printf("sum2: %v\n", sum(1))
	fmt.Printf("sum3: %v\n", sum(1, 2, 3))
}

func sum(args ...int) int {
	ret := 0
	for _, i := range args {
		ret += i
	}
	return ret
}
