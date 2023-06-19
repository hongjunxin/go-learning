package main

import "C"

import (
	"fmt"

	_ "github.com/hongjunxin/go-learning/cgo/lib/number"
)

func main() {
	println("Done")
}

//export goPrintln
func goPrintln(s *C.char) {
	fmt.Println("goPrintln:", C.GoString(s))
}
