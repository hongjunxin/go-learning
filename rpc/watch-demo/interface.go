package main

import (
	"fmt"
	"math/rand"
	"net/rpc"
	"sync"
	"time"
)

const KVStoreServiceName = "github.com/hongjunxin/go-learning/rpc.KVStoreService"

type KVStoreServiceInterface interface {
	Watch(timeoutSecond int, keyChanged *string) error
	Get(key string, value *string) error
	Set(kv [2]string, reply *struct{}) error
}

type KVStoreService struct {
	m      map[string]string
	filter map[string]func(key string)
	mu     sync.Mutex
}

func (p *KVStoreService) Watch(timeoutSecond int, keyChanged *string) error {
	id := fmt.Sprintf("watch-%s-%03d", time.Now(), rand.Int())
	ch := make(chan string, 10) // buffered

	p.mu.Lock()
	p.filter[id] = func(key string) { ch <- key }
	p.mu.Unlock()

	select {
	case <-time.After(time.Duration(timeoutSecond) * time.Second):
		return fmt.Errorf("timeout")
	case key := <-ch:
		*keyChanged = key
		return nil
	}
}

func (p *KVStoreService) Get(key string, value *string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	if v, ok := p.m[key]; ok {
		*value = v
		return nil
	}
	return fmt.Errorf("not found")
}

func (p *KVStoreService) Set(kv [2]string, reply *struct{}) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	key, value := kv[0], kv[1]
	if oldValue := p.m[key]; oldValue != value {
		for _, fn := range p.filter {
			// 当修改某个 key 对应的值时会调用每一个过滤器函数
			fn(key)
		}
	}
	p.m[key] = value
	return nil
}

type KVStoreServiceClient struct {
	Client *rpc.Client
}

func (p *KVStoreServiceClient) Watch(timeoutSecond int, keyChanged *string) error {
	return p.Client.Call(KVStoreServiceName+".Watch", timeoutSecond, keyChanged)
}

func (p *KVStoreServiceClient) Set(kv [2]string, reply *struct{}) error {
	return p.Client.Call(KVStoreServiceName+".Set", kv, reply)
}
