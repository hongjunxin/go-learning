package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hongjunxin/go-learning/trace/opentracing/common"
	"github.com/opentracing/opentracing-go"
)

func init() {
	os.Setenv("MY_PROJECT_NAME", "appA")
}

func main() {
	helloTo := ""
	if len(os.Args) != 2 {
		helloTo = "Jaeger"
	} else {
		helloTo = os.Args[1]
	}

	tracer, closer, err := common.JaegerInit()
	defer closer.Close()
	if err != nil {
		fmt.Printf("JaegerInit failed, err='%v'\n", err)
		os.Exit(-1)
	}

	span := tracer.StartSpan("main") // 设置 operation 名字，即方法名
	span.SetTag("hello-to", helloTo) // 可以设置额外的 tag
	defer span.Finish()              // 必须执行这一步
	//defer span.FinishWithOptions(opentracing.FinishOptions{FinishTime: time.Now()})

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	helloStr := methodA(ctx, helloTo)
	methodB(ctx, helloStr)
}

func methodA(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "methodA")
	defer span.Finish()

	fmt.Printf("parent span info: %v\n", common.GetSpanInfoFromContext(ctx))
	fmt.Printf("    my span info: %v\n", common.GetSpanInfo(span))

	return "Hello! " + helloTo
}

func methodB(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "methodB")
	defer span.Finish()
	fmt.Printf("[appA] methodB: %v\n", helloStr)
}
