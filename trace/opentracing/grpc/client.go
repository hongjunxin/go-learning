package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/hongjunxin/go-learning/trace/opentracing/common"
	"github.com/opentracing/opentracing-go"
)

var (
	svcaddr = flag.String("svcaddr", "localhost:9001", "server address")
)

func init() {
	os.Setenv("MY_PROJECT_NAME", "jaeger-go-gclient")
	flag.Parse()
}

func main() {
	tracer, closer, err := common.JaegerInit()
	defer closer.Close()
	if err != nil {
		log.Fatalf("JaegerInit failed, err='%v'\n", err)
	}
	for {
		span := tracer.StartSpan("main")
		ctx := opentracing.ContextWithSpan(context.Background(), span)
		_, err := common.CallSayHello(ctx, *svcaddr)
		if err != nil {
			fmt.Printf("[jaeger go] CallSayHello failed, %v", err)
			fmt.Println()
		}
		fmt.Printf("[jaeger go] client, %v\n", common.GetSpanInfoFromContext(ctx))
		span.Finish()
		time.Sleep(5 * time.Second)
	}
}
