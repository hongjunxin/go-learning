package main

import (
	context "context"
	"fmt"
	"io"
	"log"
	"time"

	grpc "google.golang.org/grpc"
)

func main() {
	done := make(chan struct{})
	go subscribeTest(done)
	time.Sleep(time.Second)
	publishTest()
	<-done
}

func publishTest() {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := NewPubsubServiceClient(conn)

	_, err = client.Publish(
		context.Background(), &String{Value: "golang: hello Go"},
	)
	if err != nil {
		log.Fatal(err)
	}
	_, err = client.Publish(
		context.Background(), &String{Value: "docker: hello Docker"},
	)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("publishTest quit")
}

func subscribeTest(done chan<- struct{}) {
	conn, err := grpc.Dial("localhost:1234", grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := NewPubsubServiceClient(conn)
	stream, err := client.Subscribe(
		context.Background(), &String{Value: "golang:"},
	)
	if err != nil {
		log.Fatal(err)
	}

	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatal(err)
		}
		fmt.Printf("client subscribe reply '%v'\n", reply.GetValue())
	}
	done <- struct{}{}
	fmt.Println("subscribeTest quit")
}
