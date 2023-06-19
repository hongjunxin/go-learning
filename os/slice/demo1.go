package main

import "fmt"

func main() {
	str := "hello"
	//str[0] = 'H' // str[0] 类型 byte，是常量，只读属性
	fmt.Printf("str[0]: %q\n", str[0])

	ss := []byte(str) // ss 会绑定新的底层数组，而不是直接用 str
	ss[0] = 'H'
	fmt.Printf("str: %v, ss: %v\n", str, string(ss)) // 输出 str: hello, ss: Hello
}
