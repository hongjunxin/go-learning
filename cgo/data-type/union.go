package main

// 推荐用 stdint.h，数据类型与 Go 兼容性好

/*
#include <stdint.h>

union B1 {
    int i;
    float f;
};

union B2 {
    int8_t i8;
    int64_t i64;
};
*/
import "C"
import (
	"fmt"
	"unsafe"
)

/*
如果需要操作 C 语言的联合类型变量，一般有三种方法：第一种是在 C 语言中定义辅助函数；
第二种是通过 Go 语言的 "encoding/binary" 手工解码成员 (需要注意大端小端问题)；
第三种是使用 unsafe 包强制转型为对应类型 (这是性能最好的方式)。

虽然 unsafe 包访问最简单、性能也最好，但是对于有嵌套联合类型的情况处理会导致问题复杂化。
对于复杂的联合类型，推荐通过在 C 语言中定义辅助函数的方式处理。
*/

func main() {
	var b1 = C.union_B1{32}
	fmt.Printf("%T\n", b1) // [4]uint8

	var b2 = C.union_B2{64}
	fmt.Printf("%T\n", b2) // [8]uint8

	fmt.Println("b1.i:", *(*C.int)(unsafe.Pointer(&b1)))
	fmt.Println("b2.i64:", *(*C.int64_t)(unsafe.Pointer(&b2)))
}
