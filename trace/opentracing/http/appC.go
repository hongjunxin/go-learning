package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/hongjunxin/go-learning/trace/opentracing/common"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

func init() {
	os.Setenv("MY_PROJECT_NAME", "appC")
}

func main() {
	tracer, closer, err := common.JaegerInit()
	defer closer.Close()
	if err != nil {
		fmt.Printf("JaegerInit failed, err='%v'\n", err)
		os.Exit(-1)
	}

	http.HandleFunc("/publish", func(w http.ResponseWriter, r *http.Request) {
		spanCtx, _ := tracer.Extract(opentracing.HTTPHeaders, opentracing.HTTPHeadersCarrier(r.Header))
		span := tracer.StartSpan("/publish", ext.RPCServerOption(spanCtx))
		defer span.Finish()

		helloStr := r.FormValue("helloStr")
		fmt.Printf("[appC] main: %v\n", helloStr)
		fmt.Printf("[appC] received header ===\n")
		common.PrinfHeaders(r.Header)
		fmt.Println()
	})

	log.Fatal(http.ListenAndServe(":8082", nil))
}
