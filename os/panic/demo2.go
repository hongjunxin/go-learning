package main

import (
	"errors"
	"fmt"
)

func main() {
	fmt.Println("Enter function main.")
	// 我们要尽量把defer语句写在函数体的开始处，因为在
	// 引发 panic 的语句之后的所有语句，都不会有任何执行机会。
	defer func() {
		fmt.Println("Enter defer function.")
		if p := recover(); p != nil {
			fmt.Printf("panic: %s\n", p)
		}

		defer func() {
			fmt.Println("Enter inner defer function")
			if p := recover(); p != nil {
				fmt.Printf("panic: %s\n", p)
			}
			fmt.Println("Exit inner defer function")
		}()
		panic(errors.New("inner painc"))

		fmt.Println("Exit defer function.")
	}()
	// 引发panic。
	panic(errors.New("something wrong"))
	fmt.Println("Exit function main.")
}
