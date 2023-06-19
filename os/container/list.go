package main

import (
	"container/list"
	"fmt"
)

func main() {
	var list list.List
	list.PushBack("a")
	e := list.PushBack("b")
	list.InsertAfter("c", e)
	d := list.InsertBefore("d", e)
	printList(&list)
	list.MoveToFront(d)
	printList(&list)
	list.Len()
}

func printList(l *list.List) {
	for i := l.Front(); i != nil; i = i.Next() {
		fmt.Printf("%v,", i.Value)
	}
	fmt.Println()
}
