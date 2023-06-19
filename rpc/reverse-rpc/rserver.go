package main

import (
	"net"
	"net/rpc"
	"time"

	t "github.com/hongjunxin/go-learning/rpc"
)

// 启动反向 rpc 服务
func main() {
	rpc.Register(new(t.HelloService))

	for {
		// 连接外网 rpc 服务
		conn, _ := net.Dial("tcp", "localhost:1234")
		if conn == nil {
			time.Sleep(time.Second)
			continue
		}
		// 作为 server 等待外网的 client 请求
		rpc.ServeConn(conn)
		conn.Close()
	}
}
