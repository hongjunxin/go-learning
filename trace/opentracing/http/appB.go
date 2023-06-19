package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"

	"github.com/hongjunxin/go-learning/trace/opentracing/common"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func init() {
	os.Setenv("MY_PROJECT_NAME", "appB")
}

func main() {
	tracer, closer, err := common.JaegerInit()
	defer closer.Close()
	if err != nil {
		fmt.Printf("JaegerInit failed, err='%v'\n", err)
		os.Exit(-1)
	}

	http.HandleFunc("/format", func(w http.ResponseWriter, r *http.Request) {
		// uber-trace-id -> span
		// x-lane-env -> span.lanEnv
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("/format", ext.RPCServerOption(spanCtx))
		fmt.Printf("appB/format get baggage item 'name': %v\n", span.BaggageItem("name"))
		defer span.Finish()

		helloTo := r.FormValue("helloTo")
		helloStr := fmt.Sprintf("Hello, %s!", helloTo)
		w.Write([]byte(helloStr))
		ctx := opentracing.ContextWithSpan(r.Context(), span)
		fmt.Printf("[appB] recieved header ===\n")
		common.PrinfHeaders(r.Header)
		callAppC(helloStr, ctx)
	})

	log.Fatal(http.ListenAndServe(":8081", nil))
}

func callAppC(helloStr string, ctx context.Context) {
	v := url.Values{}
	v.Set("helloStr", helloStr)
	url := "http://localhost:8082/publish?" + v.Encode()

	//req, err := http.NewRequest("GET", url, nil)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		fmt.Printf("[appB] http request failed, err='%v'\n", err)
	}

	span := opentracing.SpanFromContext(ctx)
	ext.SpanKindRPCClient.Set(span)
	ext.HTTPUrl.Set(span, url)
	ext.HTTPMethod.Set(span, "GET")

	// 重点!!!必须将 span 的数据 inject 到 http.NewRequest 的请求头部中
	// 这样 uber-trace-id 这个 header 才会被发送到下游服务
	span.Tracer().Inject(span.Context(), opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(req.Header))

	if _, err := common.DoRequest(req); err != nil {
		ext.LogError(span, err)
		fmt.Printf("[appB] http request failed, err='%v'\n", err)
	}
}
