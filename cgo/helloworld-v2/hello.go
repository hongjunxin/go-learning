package main

//void SayHello(char* s);
import "C"

import (
	"fmt"
)

// 这里的 C.SayHello 实际上是 main.C.SayHello
// 在Go语言中方法是依附于类型存在的，不同Go包中引入的虚拟的C包的类型却是不同的（main.C 不等于 other_package.C），
// 这导致从它们延伸出来的Go类型也是不同的类型（*main.C.char 不等 *other_package.C.char）

func main() {
	C.SayHello(C.CString("Hello, World\n"))
}

//export SayHello
func SayHello(s *C.char) {
	fmt.Print(C.GoString(s))
}
