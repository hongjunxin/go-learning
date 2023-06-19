package main

import (
	"fmt"
	"time"
)

func main() {
	str := "golang"

	go func() {
		for {
			str = fmt.Sprintf("c++ %v", time.Now())
			//fmt.Println(str)
		}
	}()
	go func() {
		for {
			fmt.Println(str)
			time.Sleep(time.Millisecond * 100)
		}
	}()
	time.Sleep(time.Second * 2)
}
