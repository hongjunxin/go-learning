package main

import (
	"fmt"
	"net/http"

	"go.elastic.co/apm/module/apmhttp"
)

// 设置多个处理器函数
func handler1(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "1 欢迎访问 www.ydook.com !")
}

func handler2(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "2 欢迎访问 www.ydook.com !")
}

func handler3(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "3 欢迎访问 www.ydook.com !")
}

func main() {
	// 设置多路复用处理函数
	mux := http.NewServeMux()

	mux.HandleFunc("/h1", handler1)
	mux.HandleFunc("/h2", handler2)
	mux.HandleFunc("/h3", handler3)

	// 设置服务器
	server := &http.Server{
		Addr:    "127.0.0.1:8000",
		Handler: mux,
	}

	apmhttp.Wrap(mux)

	// 设置服务器监听请求端口
	server.ListenAndServe()
}
