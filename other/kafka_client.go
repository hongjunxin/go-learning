package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/Shopify/sarama"
)

type ProducerWrap struct {
	p      sarama.AsyncProducer
	client sarama.Client
	sync.Mutex
	closed                              bool
	signals                             chan os.Signal
	enqueued, successes, producerErrors int
}

var (
	producer ProducerWrap
)

func init() {
	if err := producer.init(); err != nil {
		os.Exit(255)
	}
}

func main() {
	c := make(chan struct{})
	go producer.run(c)

	quit := false
	for !quit {
		message := &sarama.ProducerMessage{Topic: "topic-eg", Value: sarama.StringEncoder("testing 123")}
		sendKafKa(message)
		time.Sleep(time.Second * 3)
		select {
		case <-c:
			quit = true
		default:
		}
	}
	fmt.Println("main quit")
}

func sendKafKa(msg *sarama.ProducerMessage) {
	producer.Lock()
	defer producer.Unlock()
	if !producer.closed {
		producer.p.Input() <- msg
	}
}

func (producer *ProducerWrap) init() error {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	var err error

	//broker := []string{"10.110.223.231:9092", "10.110.208.119:9092", "10.110.194.201:9092"}
	broker := []string{"10.110.194.201:9092", "10.110.208.119:9092"}
	producer.client, err = sarama.NewClient(broker, config)
	if err != nil {
		fmt.Println("can't init  NewClient for %s, error:%s", broker, err)
		return err
	}
	producer.p, err = sarama.NewAsyncProducerFromClient(producer.client)
	if err != nil {
		return err
	}

	producer.signals = make(chan os.Signal, 1)
	signal.Notify(producer.signals, os.Interrupt)
	return nil
}

func (producer *ProducerWrap) run(c chan<- struct{}) {
	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		for range producer.p.Successes() {
			producer.successes++
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for err := range producer.p.Errors() {
			log.Println(err)
			producer.producerErrors++
		}
	}()

	quit := false
	for !quit {
		select {
		case <-producer.signals:
			producer.Lock()
			defer producer.Unlock()
			producer.p.AsyncClose()
			producer.client.Close()
			producer.closed = true
			wg.Wait()
			quit = true
		}
	}
	c <- struct{}{}
}
