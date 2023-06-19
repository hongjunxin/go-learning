package main

import (
	"log"
	"net"
	"net/rpc"
)

func NewKVStoreService() *KVStoreService {
	return &KVStoreService{
		m:      make(map[string]string),
		filter: make(map[string]func(key string)),
	}
}

func RegisterKVService(svc KVStoreServiceInterface) error {
	return rpc.RegisterName(KVStoreServiceName, svc)
}

func main() {
	RegisterKVService(NewKVStoreService())
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
