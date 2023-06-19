package main

import "fmt"

func main() {
	s1 := make([]int, 5, 8)
	for i := 0; i < len(s1); i++ {
		s1[i] = i
	}
	fmt.Printf("len(s1)=%d, cap(s1)=%d\n", len(s1), cap(s1)) // 输出 len(s1)=5, cap(s1)=8

	s2 := s1[0:3]                                            // s1 和 s2 绑定同一个底层数组
	fmt.Printf("len(s2)=%d, cap(s2)=%d\n", len(s2), cap(s2)) // 输出 len(s2)=3, cap(s2)=8
	s2[0] = 10
	fmt.Printf("s1[0]=%d, s2[0]=%d\n", s1[0], s2[0]) // 都输出 10

	s2 = append(s2, 50)
	fmt.Printf("s1[3]=%d, s2[3]=%d\n", s1[3], s2[3]) // 都输出 50

	s3 := s1[0:3]
	for i := 5; i < 10; i++ {
		// appeng 需要扩容时，返回新的切片，新切片绑定新的底层数组
		s1 = append(s1, i)
	}
	fmt.Printf("len(s1)=%d, cap(s1)=%d\n", len(s1), cap(s1))
	s3[0] = 100
	// 输出 s1[0]=10, s2[0]=100, s3[0]=100
	fmt.Printf("s1[0]=%d, s2[0]=%d, s3[0]=%d\n", s1[0], s2[0], s3[0])

	s4 := make([]int, len(s3), cap(s3))
	copy(s4, s3)
	fmt.Printf("s4: %v\n", s4) // 输出 s4: [100 1 2]
	s3[0] = 0
	fmt.Printf("s3[0]=%d, s4[0]=%d\n", s3[0], s4[0]) // 输出 s3[0]=0, s4[0]=100
}
