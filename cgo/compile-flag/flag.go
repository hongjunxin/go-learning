package main

// 在 linux 环境下，将定义 CGO_OS_LINUX=1

/*
#cgo windows CFLAGS: -DCGO_OS_WINDOWS=1
#cgo darwin CFLAGS: -DCGO_OS_DARWIN=1
#cgo linux CFLAGS: -DCGO_OS_LINUX=1

#if defined(CGO_OS_WINDOWS)
    const char* os = "windows";
#elif defined(CGO_OS_DARWIN)
    const char* os = "darwin";
#elif defined(CGO_OS_LINUX)
    const char* os = "linux";
#else
#	error(unknown os)
#endif
*/
import "C"
import "fmt"

func main() {
	fmt.Println(C.GoString(C.os))
}
