package main

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
)

func main() {
	cfg := nsq.NewConfig()

	consumer, err := nsq.NewConsumer("topic", "channel", cfg)
	if err != nil {
		log.Fatal(err)
	}

	consumer.AddHandler(nsq.HandlerFunc(func(message *nsq.Message) error {
		if len(message.Body) == 0 {
			return nil
		}

		fmt.Println("double value result: ", message.Body[0]*2)
		return nil
	}))

	if err := consumer.ConnectToNSQLookupd("127.0.0.1:4161"); err != nil {
		log.Fatal(err)
	}

	consumer.Stop()
}
