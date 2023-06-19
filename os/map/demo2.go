package main

import "fmt"

func main() {
	// channel 类型可以用作 map 的 key
	m := make(map[chan int]string)
	ch1 := make(chan int)
	ch2 := make(chan int)
	m[ch1] = "channel 1"
	m[ch2] = "channel 2"

	fmt.Printf("m[ch1]='%v', m[ch2]='%v'\n", m[ch1], m[ch2])
}
