package main

import (
	context "context"
	"fmt"
	"io"
	"log"
	"net"

	grpc "google.golang.org/grpc"
)

type HelloServiceImpl struct {
	// 推荐这里增加这个成员，增强可读性，表明应该要重写该成员的方法。
	UnimplementedHelloServiceServer
}

func (p *HelloServiceImpl) Hello(
	ctx context.Context, args *String,
) (*String, error) {
	reply := &String{Value: "hello:" + args.GetValue()}
	return reply, nil
}

func (p *HelloServiceImpl) Channel(stream HelloService_ChannelServer) error {
	fmt.Println("server Channel() called")
	for {
		args, err := stream.Recv()
		if err != nil {
			if err == io.EOF {
				fmt.Println("server received io.EOF, close channel")
				// client stream.Recv() will receive io.EOF
				// if Channel return nil
				return nil
			}
			return err
		}
		reply := &String{Value: "hello:" + args.GetValue()}
		err = stream.Send(reply)
		if err != nil {
			return err
		}
	}
}

func main() {
	grpcServer := grpc.NewServer()

	// grpcServer: 负责 grpc 通信
	// HelloServiceImpl: 负责 grpc server 服务项的实现
	RegisterHelloServiceServer(grpcServer, &HelloServiceImpl{})

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
