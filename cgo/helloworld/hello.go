package main

// import "C" 导入语句需要单独一行

//#include "hello.h"
import "C"

// 在本目录下执行 go run .

func main() {
	C.SayHello(C.CString("Hello, World\n"))
}
