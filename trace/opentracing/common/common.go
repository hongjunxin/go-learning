package common

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/hongjunxin/go-learning/trace/opentracing/jaegergrpc"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"google.golang.org/grpc"
	pb "google.golang.org/grpc/examples/helloworld/helloworld"
	"google.golang.org/grpc/metadata"
)

var (
	agentHost                     = "localhost"
	agentPort                     = "6831"
	myProjectName                 = "jaeger-client-go"
	myProjectEnvName              = "rd-default-test"
	myEnvName                     = "qa"
	jaegerSamplerType             = "remote"
	jaegerSamplerParam            = float64(1)
	jaegerSamplingServerURL       = "http://localhost/api/v1/trace/sampling"
	jaegerSamplingRefreshInterval = int64(2)
)

func jaegerInitBase() {
	if len(os.Getenv("JAEGER_AGENT_HOST")) != 0 {
		agentHost = os.Getenv("JAEGER_AGENT_HOST")
	}
	if len(os.Getenv("JAEGER_AGENT_PORT")) != 0 {
		agentPort = os.Getenv("JAEGER_AGENT_PORT")
	}
	if len(os.Getenv("MY_PROJECT_NAME")) != 0 {
		myProjectName = os.Getenv("MY_PROJECT_NAME")
	}
	if len(os.Getenv("MY_PROJECT_ENV_NAME")) != 0 {
		myProjectEnvName = os.Getenv("MY_PROJECT_ENV_NAME")
	}
	if len(os.Getenv("MY_ENV_NAME")) != 0 {
		myEnvName = os.Getenv("MY_ENV_NAME")
	}
	if len(os.Getenv("JAEGER_SAMPLER_TYPE")) != 0 {
		jaegerSamplerType = os.Getenv("JAEGER_SAMPLER_TYPE")
	}
	if len(os.Getenv("JAEGER_SAMPLER_PARAM")) != 0 {
		var err error
		jaegerSamplerParam, err = strconv.ParseFloat(os.Getenv("JAEGER_SAMPLER_PARAM"), 64)
		if err != nil {
			jaegerSamplerParam = 1
		}
	}
	if len(os.Getenv("JAEGER_SAMPLING_ENDPOINT")) != 0 {
		jaegerSamplingServerURL = os.Getenv("JAEGER_SAMPLING_ENDPOINT")
	}
	if len(os.Getenv("JAEGER_SAMPLER_REFRESH_INTERVAL")) != 0 {
		var err error
		jaegerSamplingRefreshInterval, err = strconv.ParseInt(os.Getenv("JAEGER_SAMPLER_REFRESH_INTERVAL"), 10, 64)
		if err != nil {
			jaegerSamplingRefreshInterval = 2
		}
	}
}

// JaegerInit returns an instance of Jaeger Tracer that samples 100% of traces at default.
func JaegerInit() (opentracing.Tracer, io.Closer, error) {
	jaegerInitBase()
	cfg := &config.Configuration{
		Gen128Bit:   true,
		ServiceName: myProjectName,
		Sampler: &config.SamplerConfig{
			Type:                    jaegerSamplerType,
			Param:                   jaegerSamplerParam,
			SamplingServerURL:       fmt.Sprintf("%v?service=%v&env=%v", jaegerSamplingServerURL, myProjectName, myEnvName),
			SamplingRefreshInterval: time.Duration(jaegerSamplingRefreshInterval) * time.Second,
		},
		Reporter: &config.ReporterConfig{
			LogSpans:           true,
			LocalAgentHostPort: fmt.Sprintf("%v:%v", agentHost, agentPort),
		},
		Tags: []opentracing.Tag{{Key: "env", Value: myEnvName},
			{Key: "projectEnv", Value: myProjectEnvName}},
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		return nil, nil, err
	}
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer, nil
}

func JaegerInitFromEnv() (opentracing.Tracer, io.Closer) {
	cfg, err := config.FromEnv()
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger from env: %v\n", err))
	}
	tracer, closer, err := cfg.NewTracer()
	if err != nil {
		panic(fmt.Sprintf("ERROR: cannot init Jaeger from env: %v\n", err))
	}
	// 这里顺便将 tracer 放入全局范围，opentracing 其他 api 的内部
	// 实现会使用该全局的 tracer
	opentracing.SetGlobalTracer(tracer)
	return tracer, closer
}

// Do executes an HTTP request and returns the response body.
// Any errors or non-200 status code result in an error.
func DoRequest(req *http.Request) ([]byte, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("StatusCode: %d, Body: %s", resp.StatusCode, body)
	}

	return body, nil
}

func CallSayHello(ctx context.Context, svcaddr string) (string, error) {
	conn, err := grpc.Dial(svcaddr, grpc.WithInsecure(),
		grpc.WithUnaryInterceptor(jaegergrpc.UnaryClientInterceptor()))
	if err != nil {
		return "", fmt.Errorf("connect %v failed, %v", svcaddr, err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	r, err := c.SayHello(ctx, &pb.HelloRequest{Name: "jaeger"})
	if err != nil {
		return "", fmt.Errorf("call %v SayHello failed, %v", svcaddr, err)
	}
	return r.GetMessage(), nil
}

func PrintMD(md metadata.MD) {
	for k, v := range md {
		fmt.Printf("%v: %v\n", k, v[0])
	}
	fmt.Println()
}

func PrinfHeaders(h http.Header) {
	for k, v := range h {
		fmt.Printf("%v: %v\n", k, v[0])
	}
	fmt.Println()
}

func GetSpanInfo(span opentracing.Span) string {
	spanCtx := span.Context()
	// jaeger-client-go 是 opentracing-go trace 部分的实现
	if sp, ok := spanCtx.(jaeger.SpanContext); ok {
		return fmt.Sprintf("TraceID='%v' SpanID='%v' ParentID='%v'",
			sp.TraceID().String(), sp.SpanID().String(), sp.ParentID().String())
	}
	return "print span failed"
}

func GetSpanInfoFromContext(ctx context.Context) string {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return "context lack of span"
	}
	return GetSpanInfo(span)
}
