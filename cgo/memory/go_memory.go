package main

/*
extern int* getGoPtr();

static void Main() {
    int* p = getGoPtr();
    *p = 42;
}
*/
import "C"

func main() {
	C.Main()
}

//export getGoPtr
func getGoPtr() *C.int {
	return new(C.int)
}

/*
其中 getGoPtr 返回的虽然是 C 语言类型的指针，但是内存本身是从 Go 语言的 new 函数分配，
也就是由 Go 语言运行时统一管理的内存。然后我们在 C 语言的 Main 函数中调用了 getGoPtr 函数，
此时默认将发送运行时异常：panic: runtime error: cgo result has Go pointer

除非用以下方式启动
GODEBUG=cgocheck=0 go run go_memory.go
*/
