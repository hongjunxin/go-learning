package main

import (
	"fmt"
	"log"
	"net/rpc"

	t "github.com/hongjunxin/go-learning/rpc"
)

func DialHelloService(network, address string) (*t.HelloServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &t.HelloServiceClient{Client: c}, nil
}

func main() {
	client, err := DialHelloService("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	var reply string
	err = client.Hello("hello", &reply)
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("reply='%v'\n", reply)
	}
}
