package nsq_example

import (
	"fmt"
	"github.com/nsqio/go-nsq"
	"log"
	"testing"
	"time"
)

func TestProducerPublish(t *testing.T) {
	cfg := nsq.NewConfig()

	producer, err := nsq.NewProducer("127.0.0.1:4150", cfg)
	if err != nil {
		log.Fatal(err)
	}

	topic := "topic"
	for i := 1; i <= 100; i++ {
		if err := producer.Publish(topic, []byte{byte(i)}); err != nil {
			log.Fatal(err)
		}
	}

	producer.Stop()
}

func TestConsumer(t *testing.T) {
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

	if err := consumer.ConnectToNSQD("127.0.0.1:4150"); err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 5)
	consumer.Stop()
}
