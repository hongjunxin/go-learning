package main

import (
	"log"
	"os"
	"time"
)

func main() {
	f, err := os.OpenFile("tmp.log", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		log.Fatalf("open tmp.log failed, err=%v\n", err)
	}
	log.SetOutput(f)

	go func() {
		for {
			// 已经获取了 tmp.log 的 fd，所以在持续的写入过程中，就算 tmp.log 被其他进程
			// rename 了，也不妨碍这里继续写入。或者说其他进程 rename 后，用 tail -f 查看
			// rename 后的文件时，还是能看到这里持续写入的 log
			log.Printf("time: %v\n", time.Now())
			time.Sleep(3 * time.Second)
		}
	}()

	for {
		time.Sleep(3 * time.Second)
	}
}
