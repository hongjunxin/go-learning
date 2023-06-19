package main

func main() {

}

// 返回值部分报 mixed named and unnamed parameters
// func a(a, b int) (c int, error)

func b(a, b int) (int, error) {
	return 0, nil
}

func c(a, b int, s string) error {
	return nil
}
