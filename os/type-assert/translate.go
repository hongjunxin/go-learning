package main

import "fmt"

var container = []string{"zero", "one", "two"}

func main() {
	container := map[int]string{0: "zero", 1: "one", 2: "two"}

	// 类型断言表达式的语法形式是x.(T)。其中的x代表要被判断类型的值。
	// 这个值当下的类型必须是接口类型的，不过具体是哪个接口类型其实是无所谓的。
	if value, ok := interface{}(container).(map[int]string); ok {
		fmt.Printf("translate ok, value: %v\n", value)
	}
	fmt.Printf("The element is %q.\n", container[1])
}
