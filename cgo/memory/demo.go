package main

/*
#include<stdio.h>

void printString(const char* s, int n) {
    int i;
    for(i = 0; i < n; i++) {
        putchar(s[i]);
    }
    putchar('\n');
}
*/
import "C"
import (
	"reflect"
	"unsafe"
)

/*
需要小心的是在取得 Go 内存后需要马上传入 C 语言函数，不能保存到临时变量后再间接传入 C 语言函数。
因为 CGO 只能保证在 C 函数调用之后被传入的 Go 语言内存不会发生移动，它并不能保证在传入 C 函数
之前内存不发生变化。
*/

func printString(s string) {
	p := (*reflect.StringHeader)(unsafe.Pointer(&s))
	C.printString((*C.char)(unsafe.Pointer(p.Data)), C.int(len(s)))
}

func main() {
	s := "hello"
	printString(s)
}
