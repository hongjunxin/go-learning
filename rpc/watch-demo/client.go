package main

import (
	"fmt"
	"log"
	"net/rpc"
	"strconv"
	"time"
)

func DialKVStoreService(network, address string) (*KVStoreServiceClient, error) {
	c, err := rpc.Dial(network, address)
	if err != nil {
		return nil, err
	}
	return &KVStoreServiceClient{Client: c}, nil
}

func main() {
	client, err := DialKVStoreService("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	go func() {
		var keyChanged string
		err := client.Watch(30, &keyChanged)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("watch:", keyChanged)
	}()

	if err := client.Set([2]string{"abc", strconv.Itoa(int(time.Now().Unix()))},
		new(struct{})); err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 3)
}
