package main

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// NATS_ADDRESS = nats://nats.<namespace>.svc.cluster.local:4222
	natsAddress := os.Getenv("NATS_ADDRESS")
	if natsAddress == "" {
		log.Fatal("NATS address cannot be empty")
	}
	nc, err := nats.Connect(natsAddress)
	if err != nil {
		log.Fatalf("failed to connect to NATS: %s", err)
	}
	defer nc.Close()

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("failed to create JS context: %s", err)
	}

	sub, err := js.PullSubscribe("ORDERS.*", "PULL_CONSUMER")
	if err != nil {
		log.Fatalf("failed to create pull subscription: %s", err)
	}

	batch := 5
	for i := 0; i < 50; i++ {
		msgs, err := sub.Fetch(batch, nats.MaxWait(2*time.Second))
		if err != nil {
			log.Fatalf("failed to get message: %s", err)
		}
		// Ack messages.
		for _, msg := range msgs {
			msg.AckSync()
		}
		time.Sleep(1 * time.Second)
	}
}
