package main

import (
	"fmt"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
	"github.com/hongjunxin/go-learning/net/netx"
)

func initTracerProvider() (*tracesdk.TracerProvider, error) {
	return tracerProviderToAgent("localhost", "6831")
}

// 上报给 jaeger agent
func tracerProviderToAgent(agentHost, agentPort string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithAgentEndpoint(jaeger.WithAgentHost(agentHost), jaeger.WithAgentPort(agentPort)))
	if err != nil {
		return nil, err
	}
	hostIP := ""
	if ip, err := netx.HostIP(); err != nil {
		hostIP = ip.String()
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(os.Getenv("MY_PROJECT_NAME")),
			attribute.String("projectEnv", os.Getenv("MY_PROJECT_ENV_NAME")),
			attribute.String("env", os.Getenv("MY_ENV_NAME")),
			attribute.String("ip", hostIP),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetErrorHandler(otel.ErrorHandlerFunc(func(err error) {
		fmt.Printf("[otel] error: %v\n", err)
	}))
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}

// 上报给 jaeger collector
func TracerProviderToCollector(url string) (*tracesdk.TracerProvider, error) {
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String(os.Getenv("MY_PROJECT_NAME")),
			attribute.String("projectEnv", os.Getenv("MY_PROJECT_ENV_NAME")),
			attribute.String("env", os.Getenv("MY_ENV_NAME")),
		)),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp, nil
}
