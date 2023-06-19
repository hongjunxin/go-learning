package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	zipkin "github.com/openzipkin/zipkin-go"
	zipkinhttp "github.com/openzipkin/zipkin-go/middleware/http"
	"github.com/openzipkin/zipkin-go/model"
	"github.com/openzipkin/zipkin-go/propagation/b3"
	"github.com/openzipkin/zipkin-go/reporter"
	httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
)

func main() {
	clientDemo()
}

func doSomeWork(context.Context) {}

func clientDemo() {
	//reporter := logreporter.NewReporter(log.New(os.Stderr, "", log.LstdFlags))
	// httpreporter "github.com/openzipkin/zipkin-go/reporter/http"
	reporter := httpreporter.NewReporter("http://localhost:9411/api/v2/spans") // report to jaeger-collector zipkin-server
	defer reporter.Close()
	os.Getenv("projectEnv")

	// set-up the remote endpoint for the service we're about to call
	endpoint, err := zipkin.NewEndpoint("gotrace_web", "127.0.0.1:8686")
	if err != nil {
		log.Fatalf("unable to create remote endpoint: %+v\n", err)
	}
	sampler, err := zipkin.NewBoundarySampler(1, time.Now().UnixNano())
	if err != nil {
		log.Fatalf("unable to create sampler: %+v\n", err)
	}
	tracer, err := zipkin.NewTracer(
		reporter,
		zipkin.WithLocalEndpoint(endpoint),
		zipkin.WithSampler(sampler),
	)
	if err != nil {
		log.Fatalf("unable to create tracer: %+v\n", err)
	}

	client, err := zipkinhttp.NewClient(tracer, zipkinhttp.ClientTrace(true))
	if err != nil {
		log.Fatalf("unable to create client: %+v\n", err)
	}
	req, err := http.NewRequest("GET", "http://127.0.0.1:8686/span4", nil)
	if err != nil {
		log.Fatalf("unable to create http request: %+v\n", err)
	}
	res, err := client.DoWithAppSpan(req, "")
	if err != nil {
		log.Fatalf("unable to do http request: %+v\n", err)
	}
	res.Body.Close()
}

func ExampleTracerOption() {
	tracer, _ := zipkin.NewTracer(
		reporter.NewNoopReporter(),
		zipkin.WithNoopSpan(true),
	)

	span := tracer.StartSpan("some_operation")
	var client http.Client
	req, err := http.NewRequest("GET", "http://127.0.0.1:8686/span4", nil)
	injector := b3.InjectHTTP(req)
	injector(span.Context())
	if err != nil {
		log.Fatalf("unable to create http request: %+v\n", err)
	}
	res, err := client.Do(req)
	if err != nil {
		log.Fatalf("unable to do http request: %+v\n", err)
	}
	res.Body.Close()
	span.Finish()
}

func ExampleNewContext() {
	var (
		tracer, _ = zipkin.NewTracer(reporter.NewNoopReporter())
		ctx       = context.Background()
	)

	// span for this function
	span := tracer.StartSpan("ExampleNewContext")
	defer span.Finish()

	// add span to Context
	ctx = zipkin.NewContext(ctx, span)

	// pass along Context which holds the span to another function
	doSomeWork(ctx)

	// Output:
}

func ExampleSpanOption() {
	tracer, _ := zipkin.NewTracer(reporter.NewNoopReporter())

	// set-up the remote endpoint for the service we're about to call
	endpoint, err := zipkin.NewEndpoint("otherService", "localhost:80")
	if err != nil {
		log.Fatalf("unable to create remote endpoint: %+v\n", err)
	}

	// start a client side RPC span and use RemoteEndpoint SpanOption
	span := tracer.StartSpan(
		"some-operation",
		zipkin.RemoteEndpoint(endpoint),
		zipkin.Kind(model.Client),
	)
	// ... call other service ...
	span.Finish()

	// Output:
}
