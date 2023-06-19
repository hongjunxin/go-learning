package main

import (
	"log"
	"net"
	"net/rpc"
	"net/rpc/jsonrpc"

	t "github.com/hongjunxin/go-learning/rpc"
)

func main() {
	rpc.RegisterName("HelloService", new(t.HelloService))
	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go rpc.ServeCodec(jsonrpc.NewServerCodec(conn))
	}
}
