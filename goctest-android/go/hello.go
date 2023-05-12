package main

import "C"
import "fmt"

func main() {}

//export SayHi
func SayHi() {
	fmt.Println("Hi, I'am golang")
}
