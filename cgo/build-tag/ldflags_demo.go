/*
-X importpath.name=value 编译期设置变量的值
-s disable symbol table 禁用符号表
-w disable DWARF generation 禁用调试信息
*/

// go build -ldflags "-s -w -X 'main.BUILD_TIME=1900.1.1' -X 'main.GO_VERSION=v0.0.1' -X 'main.VERSION=1.0.0'"  -o t ldflags_demo.go

package main

import "fmt"

var (
	VERSION    string
	BUILD_TIME string
	GO_VERSION string
)

func main() {
	fmt.Printf("%s\n%s\n%s\n", VERSION, BUILD_TIME, GO_VERSION)
}
