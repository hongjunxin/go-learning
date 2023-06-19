package main

import (
	context "context"
	"log"
	"net"
	"strings"
	"time"

	"github.com/moby/pubsub"
	"google.golang.org/grpc"
)

type PubsubService struct {
	pub *pubsub.Publisher
	m   map[string]chan interface{}
}

func NewPubsubService() *PubsubService {
	return &PubsubService{
		pub: pubsub.NewPublisher(100*time.Millisecond, 10),
	}
}

func (p *PubsubService) Publish(
	ctx context.Context, arg *String,
) (*String, error) {
	p.pub.Publish(arg.GetValue())
	return &String{}, nil
}

func (p *PubsubService) Subscribe(
	arg *String, stream PubsubService_SubscribeServer,
) error {
	ch := p.pub.SubscribeTopic(func(v interface{}) bool {
		if key, ok := v.(string); ok {
			if strings.HasPrefix(key, arg.GetValue()) {
				return true
			}
		}
		return false
	})
	// todo: save ch to p.m
	for v := range ch {
		if err := stream.Send(&String{Value: v.(string)}); err != nil {
			return err
		}
	}

	return nil
}

// todo:
func (p *PubsubService) CancelSubscribe() {

}

func main() {
	grpcServer := grpc.NewServer()
	RegisterPubsubServiceServer(grpcServer, NewPubsubService())

	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatal(err)
	}
	grpcServer.Serve(lis)
}
