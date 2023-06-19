package main

import (
	"fmt"
	"reflect"
)

type animal interface {
	eat()
}

type bird struct {
	name string
}

func (b *bird) eat() {
	fmt.Println("bird do eat")
}

func main() {
	b := &bird{name: "nightingale"}
	var aifa animal
	aifa = b
	t := reflect.TypeOf(aifa)
	// Comparable 是否能进行判等操作
	fmt.Printf("dynamic type: %v, comparable: %v\n", t, t.Comparable())
}
