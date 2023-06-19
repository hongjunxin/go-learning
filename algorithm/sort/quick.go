package main

import "fmt"

func main() {
	arr := []int{3, 2, 4, 0, 1, 5, 7, 8, 9, 6}
	quickSort(0, len(arr)-1, arr)
	fmt.Println(arr)
}

func doPart(begin, end int, array []int) int {
	small := begin - 1
	pivot := array[end]
	for i := begin; i < end; i++ {
		if array[i] < pivot {
			small++
			swap(small, i, array)
		} else if array[i] == pivot {
			small++
		}
	}
	small++
	swap(small, end, array)
	return small
}

func quickSort(begin, end int, array []int) {
	if begin >= end {
		return
	}
	donePoint := doPart(begin, end, array)
	quickSort(begin, donePoint-1, array)
	quickSort(donePoint+1, end, array)
}

func swap(a, b int, array []int) {
	tmp := array[a]
	array[a] = array[b]
	array[b] = tmp
}
