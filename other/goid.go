package main

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"
)

func main() {
	go func() {
		fmt.Printf("go1 id=%v\n", GetGoid())
	}()
	go func() {
		fmt.Printf("go1 id=%v\n", GetGoid())
	}()
	time.Sleep(time.Second)
}

func GetGoid() int64 {
	var (
		buf [64]byte
		n   = runtime.Stack(buf[:], false)
		stk = strings.TrimPrefix(string(buf[:n]), "goroutine ")
	)

	idField := strings.Fields(stk)[0]
	id, err := strconv.Atoi(idField)
	if err != nil {
		panic(fmt.Errorf("can not get goroutine id: %v", err))
	}

	return int64(id)
}
