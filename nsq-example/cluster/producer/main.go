package main

import (
	"github.com/nsqio/go-nsq"
	"log"
	"time"
)

func main() {
	cfg := nsq.NewConfig()

	producer, err := nsq.NewProducer("127.0.0.1:4150", cfg)
	if err != nil {
		log.Fatal(err)
	}

	topic := "topic"
	for i := 1; i <= 100; i++ {
		time.Sleep(1 * time.Second)
		if err := producer.Publish(topic, []byte{byte(i)}); err != nil {
			log.Fatal(err)
		}
	}

	producer.Stop()
}
