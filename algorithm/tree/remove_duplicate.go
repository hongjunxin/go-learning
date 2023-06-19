package main

import (
	"fmt"
)

type Node struct {
	Value int
	Next  *Node
}

func main() {
	p1 := &Node{1, nil}
	p2 := &Node{3, nil}
	p3 := &Node{3, nil}
	p4 := &Node{3, nil}
	p5 := &Node{4, nil}

	p1.Next = p2
	p2.Next = p3
	p3.Next = p4
	p4.Next = p5

	solution(p1)
	p := p1
	for p != nil {
		fmt.Printf("%v ", p.Value)
		p = p.Next
	}
	fmt.Println()
}

func solution(head *Node) *Node {
	if head == nil {
		return nil
	}
	pa := head
	pb := pa.Next
	for pb != nil {
		if pa.Value != pb.Value {
			pa = pb
			pb = pb.Next
		} else {
			pa.Next = pb.Next
			pb = pb.Next
		}
	}
	return head
}
