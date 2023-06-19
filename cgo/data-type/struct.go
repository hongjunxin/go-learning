package main

/*
在 Go 语言中，我们可以通过 C.struct_xxx 来访问 C 语言中定义的 struct xxx 结构体类型。
结构体的内存布局按照 C 语言的通用对齐规则，在 32 位 Go 语言环境 C 语言结构体也按照 32 位对齐规则，
在 64 位 Go 语言环境按照 64 位的对齐规则。对于指定了特殊对齐规则的结构体，无法在 CGO 中访问。
*/

/*
struct A {
    int i;
    float f;
	int type; // type 是 Go 语言的关键字
	int   size: 10; // 位字段无法访问
};
*/
import "C"
import "fmt"

func main() {
	var a C.struct_A
	fmt.Println(a.i)
	fmt.Println(a.f)
	// 如果结构体的成员名字中碰巧是 Go 语言的关键字，可以通过在成员名开头添加下划线来访问
	fmt.Println(a._type)
}
