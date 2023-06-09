package main

import (
	"log"
	"net"
	"net/rpc"

	t "github.com/hongjunxin/go-learning/rpc"
)

func RegisterHelloService(svc t.HelloServiceInterface) error {
	return rpc.RegisterName(t.HelloServiceName, svc)
}

func main() {
	RegisterHelloService(new(t.HelloService))

	listener, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal("ListenTCP error:", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal("Accept error:", err)
		}

		go rpc.ServeConn(conn)
	}
}
