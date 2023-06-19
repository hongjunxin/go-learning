package main

import (
	"errors"
	"fmt"
)

func main() {
	defer func() {
		if p := recover(); p != nil {
			fmt.Printf("main() caught painc: %s\n", p)
		}
	}()
	demo()
}

func demo() {
	fmt.Println("Enter function demo.")
	// 我们要尽量把defer语句写在函数体的开始处，因为在
	// 引发 panic 的语句之后的所有语句，都不会有任何执行机会。
	defer func() {
		fmt.Println("Enter defer function.")
		if p := recover(); p != nil {
			fmt.Printf("panic: %s\n", p)
			if p != errors.New("can handle") {
				// 应用场景：不是该函数能处理的 panic，则继续抛出让上层函数处理
				panic(p)
			}
		}
		fmt.Println("Exit defer function.")
	}()
	// 引发 panic
	panic(errors.New("something wrong"))
	fmt.Println("Exit function demo.")
}
