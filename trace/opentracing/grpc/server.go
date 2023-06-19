package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/hongjunxin/go-learning/trace/opentracing/common"
	"github.com/hongjunxin/go-learning/trace/opentracing/jaegergrpc"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

var (
	port = flag.Int("port", 9001, "server port")
)

func init() {
	os.Setenv("MY_PROJECT_NAME", "jaeger-go-gserver")
	flag.Parse()
}

func main() {
	_, closer, err := common.JaegerInit()
	defer closer.Close()
	if err != nil {
		fmt.Printf("JaegerInit failed, err='%v'\n", err)
		os.Exit(-1)
	}

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(jaegergrpc.UnaryServerInterceptor()))
	reflection.Register(s)
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("=== [jaeger go] SayHello server received headers ===")
		common.PrintMD(md)
	} else {
		fmt.Println("SayHello(): metadata.FromIncomingContext(ctx) failed")
	}

	_, err := common.CallSayHello(ctx, "localhost:9002")
	if err != nil {
		fmt.Printf("[jaeger go] CallSayHello failed, %v", err)
		fmt.Println()
	}
	return &pb.HelloReply{Message: "Hello: " + in.GetName()}, nil
}
