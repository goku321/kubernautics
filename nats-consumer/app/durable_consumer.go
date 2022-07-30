package main

import (
	"log"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	nc, err := nats.Connect("nats://nats.nats.svc.cluster.local:4222")
	if err != nil {
		log.Fatalf("failed to connect to NATS: %s", err)
	}

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("failed to create JS context: %s", err)
	}

	_, err = js.AddConsumer("mystream", &nats.ConsumerConfig{
		Durable: "PULL_CONSUMER",
	})
	if err != nil {
		log.Fatalf("failed to create consumer: %s\n", err)
	}

	sub, err := js.PullSubscribe("ORDERS.*", "PULL_CONSUMER")
	if err != nil {
		log.Fatalf("failed to create pull subscription: %s", err)
	}

	for i := 0; i < 1000; i++ {
		msg, err := sub.NextMsg(2 * time.Second)
		if err != nil {
			log.Fatalf("failed to get message: %s", err)
		}
		msg.Ack()
	}
}