//go:build debug
// +build debug

package main

import "fmt"

var buildMode = "debug"

func main() {
	fmt.Println("build tag demo")
}

/*
go build -tags="debug"
追加了 -tags="debug" 时才会编译这个文件

当有多个 build tag 时，我们将多个标志通过逻辑操作的规则来组合使用。

+build linux,386 darwin,!cgo
其中 linux,386 中 linux 和 386 用逗号链接表示 AND 的意思；而 linux,386 和 darwin,!cgo 之间通过空白分割来表示 OR 的意思。
*/
