package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/hongjunxin/go-learning/trace/opentracing/common"
	"github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
	"github.com/opentracing/opentracing-go/log"
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
	span.SetBaggageItem("name", "jaxon")
	defer span.Finish() // 必须执行这一步
	fmt.Printf("[appA] main: %v\n", common.GetSpanInfo(span))

	ctx := opentracing.ContextWithSpan(context.Background(), span)

	helloStr := methodA(ctx, helloTo)
	methodB(ctx, helloStr)
}

func methodA(ctx context.Context, helloTo string) string {
	span, _ := opentracing.StartSpanFromContext(ctx, "methodA")
	fmt.Printf("methodA get baggage item 'name': %v\n", span.BaggageItem("name"))
	defer span.Finish()

	v := url.Values{}
	v.Set("helloTo", helloTo)
	url := "http://localhost:8081/format?" + v.Encode()
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		panic(err.Error())
	}
	req.Header.Add("x-lane-env", "sbu-test-1")

	ext.SpanKindRPCClient.Set(span) // 设置tag "span.kind: client"
	ext.HTTPUrl.Set(span, url)      // 设置tag "http.url: http://localhost:8081/format?helloTo=x"
	ext.HTTPMethod.Set(span, "GET") // 设置tag "http.method: GET"

	// 重点!!!必须将 span 的数据 inject 到 http.NewRequest 的请求头部中
	// 这样 uber-trace-id 这个 header 才会被发送到下游服务
	span.Tracer().Inject(
		span.Context(),
		opentracing.HTTPHeaders,
		opentracing.HTTPHeadersCarrier(req.Header),
	)

	resp, err := common.DoRequest(req)
	if err != nil {
		ext.LogError(span, err)
		fmt.Printf("[appA] http request failed, err='%v'\n", err)
	}

	helloStr := string(resp)
	// 非必须项。可以将日志也上报给 jaeger agent
	span.LogFields(
		log.String("event", "string-format"),
		log.String("value", helloStr),
	)

	return helloStr
}

func methodB(ctx context.Context, helloStr string) {
	span, _ := opentracing.StartSpanFromContext(ctx, "methodB")
	defer span.Finish()
	fmt.Printf("[appA] methodB: %v\n", helloStr)
}
