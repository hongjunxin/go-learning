package rpc

import "net/rpc"

const HelloServiceName = "github.com/hongjunxin/go-learning/rpc.HelloService"

type HelloServiceInterface interface {
	Hello(request string, reply *string) error
}

type HelloService struct{}

func (p *HelloService) Hello(request string, reply *string) error {
	*reply = "hello:" + request
	return nil
}

type HelloServiceClient struct {
	Client *rpc.Client
}

func (p *HelloServiceClient) Hello(request string, reply *string) error {
	return p.Client.Call(HelloServiceName+".Hello", request, reply)
}

// 用于检查 HelloServiceClient 是否实现了 HelloServiceInterface
var _ HelloServiceInterface = (*HelloServiceClient)(nil)
