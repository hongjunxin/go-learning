package main

import (
	context "context"
	"fmt"
	"io"
	"log"
	"time"

	grpc "google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	conn, err := grpc.Dial("localhost:1234", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := NewHelloServiceClient(conn)
	reply, err := client.Hello(context.Background(), &String{Value: "hello"})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(reply.GetValue())

	streamTest(&client)
}

func streamTest(client *HelloServiceClient) {
	fmt.Println("client Channel() called")
	stream, err := (*client).Channel(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	go func() {
		for i := 0; i < 3; i++ {
			if err := stream.Send(&String{Value: "hi"}); err != nil {
				log.Fatal(err)
			}
			time.Sleep(time.Second)
		}
		stream.CloseSend()
	}()
	for {
		reply, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("client received io.EOF, quit")
				break
			}
			log.Fatal(err)
		}
		fmt.Println(reply.GetValue())
	}
}
