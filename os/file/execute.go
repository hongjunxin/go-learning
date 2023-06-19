package main

import (
	"fmt"
	"os"
	"path/filepath"
	// 不允许访问 internl 包，编译不通过
	//"github.com/hongjunxin/go-learning/os/flag/internal"
)

func main() {
	dir, base := getDirAndBaseName()
	fmt.Printf("dir: %v, basename: %v\n", dir, base)
}

func getDirAndBaseName() (string, string) {
	// 获取程序运行时的绝对路径
	executable, err := os.Executable()
	if err != nil {
		return "", ""
	}
	//internal.Hello()
	return filepath.Dir(executable), filepath.Base(executable)
}
