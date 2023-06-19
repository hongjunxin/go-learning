package main

import (
	"fmt"
	"strings"
)

func main() {
	s := "a,b,c,d"
	// 参数 n 表示最多切分出几个子串，超出的部分将不再切分。
	// 如果 n 为 0 则返回 nil, 如果 n 小于 0，则不限制切分个数，全部切分
	// 带 After 的话, 则输出的结果中会带上分隔符
	ss := strings.SplitN(s, ",", 2)
	ssn := strings.SplitAfterN(s, ",", 2)
	fmt.Println(ss)  // 输出 [a b,c,d]
	fmt.Println(ssn) // 输出 [a, b,c,d]
}
