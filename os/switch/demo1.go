package main

import "fmt"

// switch 语句会进行有限的类型转换

func main() {
	value1 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	// 1 + 3 得到的默认数据类型是 int，而 case 表达式中的数据类型是 int8
	// 所以这个 switch 不能通过编译
	switch 1 + 3 {
	case value1[0], value1[1]:
		fmt.Println("0 or 1")
	case value1[2], value1[3]:
		fmt.Println("2 or 3")
	case value1[4], value1[5], value1[6]:
		fmt.Println("4 or 5 or 6")
	}

	value2 := [...]int8{0, 1, 2, 3, 4, 5, 6}
	switch value2[4] {
	/* 如果 case 表达式中子表达式的结果值是无类型的常量，那么它的类型会被自动地转换为 switch 表达式的结果类型，
	   又由于这里 case 的整数都可以被转换为 int8 类型的值，所以这个 switch 可以通过编译。*/
	case 0, 1:
		fmt.Println("0 or 1")
	case 2, 3:
		fmt.Println("2 or 3")
	case 4, 5, 6:
		fmt.Println("4 or 5 or 6")
	}

	// type byte = uint8
	value6 := interface{}(byte(127))
	switch t := value6.(type) {
	case uint8, uint16:
		fmt.Println("uint8 or uint16")
	case byte:
		fmt.Printf("byte")
	default:
		fmt.Printf("unsupported type: %T", t)
	}
}
