// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Command jaeger is an example program that creates spans
// and uploads to Jaeger.
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/uber/jaeger-client-go/utils"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.10.0"
	"go.opentelemetry.io/otel/trace"
)

var (
	jaegerAgentHost = "localhost"
	jaegerAgentPort = "6831"
)

func init() {
	if len(os.Getenv("JAEGER_AGENT_HOST")) != 0 {
		jaegerAgentHost = os.Getenv("JAEGER_AGENT_HOST")
	}
	if len(os.Getenv("JAEGER_AGENT_PORT")) != 0 {
		jaegerAgentPort = os.Getenv("JAEGER_AGENT_PORT")
	}
	os.Setenv("MY_ENV_NAME", "qa")
	os.Setenv("MY_PROJECT_ENV_NAME", "rd-test-1")
	os.Setenv("MY_PROJECT_NAME", "otel-go-one-process")
}

func TracerProviderToCollector(url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(os.Getenv("MY_PROJECT_NAME")),
			attribute.String("projectEnv", os.Getenv("MY_PROJECT_ENV_NAME")),
			attribute.String("env", os.Getenv("MY_ENV_NAME")),
		)),
	)
	return tp, nil
}

func tracerProviderToAgent(agentHost, agentPort string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(agentHost),
		jaeger.WithAgentPort(agentPort)))
	if err != nil {
		return nil, err
	}
	hostIP := ""
	if netIP, err := utils.HostIP(); err == nil {
		hostIP = netIP.String()
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in a Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(os.Getenv("MY_PROJECT_NAME")),
			attribute.String("projectEnv", os.Getenv("MY_PROJECT_ENV_NAME")),
			attribute.String("env", os.Getenv("MY_ENV_NAME")),
			attribute.String("ip", hostIP),
		)),
		tracesdk.WithSampler(tracesdk.ParentBased(tracesdk.TraceIDRatioBased(1))),
	)

	return tp, nil
}

// func getLocalIP() (string, error) {
// 	jaegerAddr := fmt.Sprintf("%v:%v", jaegerAgentHost, jaegerAgentPort)
// 	addr, err := net.ResolveUDPAddr("udp", jaegerAddr)
// 	if err != nil {
// 		return "", err
// 	}
// 	socket, err := net.DialUDP("udp", nil, addr)
// 	if err != nil {
// 		return "", err
// 	}
// 	localAddr := socket.LocalAddr().String()
// 	var ss []string
// 	if localAddr[0] == '[' {
// 		localAddr = localAddr[1:]
// 		ss = strings.SplitN(localAddr, "]", 2)
// 	} else {
// 		ss = strings.SplitN(localAddr, ":", 2)
// 	}
// 	socket.Close()
// 	return ss[0], nil
// }

func main() {
	//tp, err := TracerProviderToCollector(os.Getenv("JAEGER_COLLECTOR_ENDPOINT"))
	tp, err := tracerProviderToAgent(jaegerAgentHost, jaegerAgentPort)
	if err != nil {
		log.Fatal(err)
	}

	// 将 TracerProvider 注册为全局变量，因为 otel 的其他实现将会
	// 在内部使用它。
	otel.SetTracerProvider(tp)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Cleanly shutdown and flush telemetry when the application exits.
	defer func(ctx context.Context) {
		// 设置超时防止应用被 tp.Shutdown() hang 住
		ctx, cancel = context.WithTimeout(ctx, time.Second*5)
		defer cancel()
		if err := tp.Shutdown(ctx); err != nil {
			log.Fatal(err)
		}
	}(ctx)

	tr := tp.Tracer("component-main") // 参数为 operation 的名字

	// 如果 ctx 中没有包含 span，则返回的 span 为 root span，如果有包含
	// 则返回的是 child span。这里返回的是 root span
	// "foo" 是该 span 的名字
	ctx, span := tr.Start(ctx, "foo")
	fmt.Printf("foo trace info: %v\n", GetSpanInfo(span))
	defer span.End()

	bar(ctx)
}

func bar(ctx context.Context) {
	// 使用全局的 TracerProvider，与 otel.SetTracerProvider() 对应。
	tr := otel.Tracer("component-bar")
	_, span := tr.Start(ctx, "bar")
	span.SetAttributes(attribute.Key("testset").String("value"))
	fmt.Printf("bar trace info: %v\n", GetSpanInfo(span))
	defer span.End()

	// Do bar...
}

func GetSpanInfo(span trace.Span) string {
	return fmt.Sprintf("TraceID='%v' SpanID='%v'",
		span.SpanContext().TraceID().String(), span.SpanContext().SpanID().String())
}
