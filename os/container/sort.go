package main

import (
	"fmt"
	"sort"
)

type IntSlice []int

func (s IntSlice) Len() int           { return len(s) }
func (s IntSlice) Less(i, j int) bool { return s[i] < s[j] }
func (s IntSlice) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func main() {
	s := []int{4, 5, 1, 7, 2, 9}
	sort.Sort(IntSlice(s))
	fmt.Println(s)
}
