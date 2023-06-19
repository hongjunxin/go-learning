package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/otel/baggage"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
)

var (
	svcaddr = flag.String("svcaddr", "localhost:9002", "server address")
)

func init() {
	flag.Parse()
}

func main() {
	os.Setenv("MY_ENV_NAME", "qa")
	os.Setenv("MY_PROJECT_ENV_NAME", "rd-test-1")
	os.Setenv("MY_PROJECT_NAME", "otel-go-gclient")
	tp, err := initTracerProvider()
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if err := tp.Shutdown(context.Background()); err != nil {
			log.Printf("Error shutting down tracer provider: %v", err)
		}
	}()

	var conn *grpc.ClientConn
	conn, err = grpc.Dial(*svcaddr, grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithUnaryInterceptor(otelgrpc.UnaryClientInterceptor()),
	)
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	bag, _ := baggage.Parse("x-lane-env=rd-test-1")
	pKey1, _ := baggage.NewKeyProperty("pkey1")
	pKey2, _ := baggage.NewKeyValueProperty("pKey2", "pValue2")
	member, _ := baggage.NewMember("x-test", "value", pKey1, pKey2)
	bag, _ = bag.SetMember(member)
	ctx = baggage.ContextWithBaggage(ctx, bag)

	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "otel"})
	if err != nil {
		log.Printf("[client]could not greet: %v\n", err)
	}
	fmt.Printf("[client] got ret '%v'\n", r.GetMessage())
}
