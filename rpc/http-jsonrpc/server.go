package main

import (
	"io"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"

	t "github.com/hongjunxin/go-learning/rpc"
)

func main() {
	rpc.RegisterName("HelloService", new(t.HelloService))

	http.HandleFunc("/jsonrpc", func(w http.ResponseWriter, r *http.Request) {
		var conn io.ReadWriteCloser = struct {
			io.Writer
			io.ReadCloser
		}{
			ReadCloser: r.Body,
			Writer:     w,
		}

		rpc.ServeRequest(jsonrpc.NewServerCodec(conn))
	})

	http.ListenAndServe(":1234", nil)
}
