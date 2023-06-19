package main

import (
	"context"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptrace"
	"os"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/net/http/httptrace/otelhttptrace"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/reflection"
)

var (
	port = flag.Int("port", 9002, "server port")
)

func init() {
	flag.Parse()
}

// server is used to implement helloworld.GreeterServer.
type server struct {
	pb.UnimplementedGreeterServer
}

func main() {
	os.Setenv("MY_ENV_NAME", "qa")
	os.Setenv("MY_PROJECT_ENV_NAME", "rd-test-1")
	os.Setenv("MY_PROJECT_NAME", "otel-go-gserver")
	tp, err := initTracerProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", fmt.Sprintf(":%v", *port))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(grpc.UnaryInterceptor(otelgrpc.UnaryServerInterceptor()))
	reflection.Register(s)
	pb.RegisterGreeterServer(s, &server{})
	log.Printf("server listening at %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// SayHello implements helloworld.GreeterServer
func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	bag := baggage.FromContext(ctx)
	for i, member := range bag.Members() {
		fmt.Printf("member[%v]: %v\n", i, member)
	}
	fmt.Println()
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		fmt.Println("=== [otel go] SayHello server received headers ===")
		printMD(md)
	} else {
		fmt.Println("SayHello(): metadata.FromIncomingContext(ctx) failed")
	}
	//doHttpRequest(ctx)
	return &pb.HelloReply{Message: "Hello: " + in.GetName()}, nil
}

func printMD(md metadata.MD) {
	for k, v := range md {
		fmt.Printf("%v: %v\n", k, v[0])
	}
	fmt.Println()
}

func doHttpRequest(ctx context.Context) {
	client := http.Client{Transport: otelhttp.NewTransport(http.DefaultTransport)}
	url := "http://localhost:8686/span3"
	tr := otel.Tracer("example/client")
	err := func(ctx context.Context) error {
		ctx, span := tr.Start(ctx, "say-hello", trace.WithAttributes(semconv.PeerServiceKey.String("ExampleService")))
		defer span.End()

		ctx = httptrace.WithClientTrace(ctx, otelhttptrace.NewClientTrace(ctx))
		req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
		res, err := client.Do(req)
		if err != nil {
			fmt.Printf("[otel go] grpc_server call '%v' failed, err='%v'\n", url, err)
		}
		_, err = ioutil.ReadAll(res.Body)
		_ = res.Body.Close()

		return err
	}(ctx)

	if err != nil {
		fmt.Printf("[otel go] grpc_server call '%v' failed, err='%v'\n", url, err)
	}
}
