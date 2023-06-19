package main

import "fmt"

func main() {
	fmt.Println("test1")
	test1()
	fmt.Println("\ntest2")
	test2()
}

func test1() {
	// 有等号的写法，mystring 等同于 string，只是名字不一样
	type mystring = string
	var (
		str1 string   = "Hi"
		str2 mystring = "Hi"
	)
	// 输出 str1 type: string, str2 tyep: string
	fmt.Printf("str1 type: %T, str2 tyep: %T\n", str1, str2)
	fmt.Printf("str1 equals str2: %v\n", str1 == str2) // 输出 true

	arry1 := []string{}
	arry2 := []mystring{"hi"}

	arry1 = arry2 // 可以直接赋值，不用进行类型转换
	fmt.Printf("arry1: %v\n", arry1)
}

func test2() {
	// 不加等号的写法，可以理解为创造了新的类型，类型不一样了，但是彼此之间可以进行类型转换
	type mystring string
	var (
		str1 string   = "Hi"
		str2 mystring = "Hello"
	)

	// 输出 str1 type: string, str2 type: main.mystring
	fmt.Printf("str1 type: %T, str2 type: %T\n", str1, str2)

	// 无法通过编译
	// mismatched types string and mystring
	//fmt.Printf("equal: %v\n", str1 == str2)

	str1 = string(str2) // 但可以进行类型转换
	fmt.Printf("str1: %v\n", str1)

	//arry1 := []string{}
	//arry2 := []mystring{"hi"}
	// 无法通过编译
	// cannot convert arry2 (variable of type []mystring) to []string
	//arry1 = []string(arry2)
}
