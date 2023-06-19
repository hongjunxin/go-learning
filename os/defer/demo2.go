package main

import "fmt"

func main() {
	fmt.Printf("out1(1): %v\n", out1(1)) // 输出 2
	fmt.Printf("out2(1): %v\n", out2(1)) // 输出 11
	fmt.Printf("out3(1): %v\n", out3(1)) // 输出 4
}

// defer 在 return 语句之后执行

func out1(a int) (ret int) {
	ret = 1
	defer func() {
		// defer 中会修改 ret 的值
		ret += a
	}()
	// 这里 return 没有对 ret 进行修改
	return
}

func out2(a int) (ret int) {
	ret = 1
	defer func() {
		// defer 又在 10 的基础上加上 a
		ret += a
	}()
	// return 将 ret 设置为 10
	return 10
}

func out3(a int) int {
	ret := 4
	defer func() {
		ret += a
	}()
	// 返回值没有指定变量名，所以返回值直接就是 return 返回的值
	return ret
}
