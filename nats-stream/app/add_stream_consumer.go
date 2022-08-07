package main

import (
	"log"
	"os"
	"time"

	"github.com/nats-io/nats.go"
)

func main() {
	// NATS_ADDRESS = nats://nats.<namespace>.svc.cluster.local:4222
	// nats.nats-jetstream-test-nats-ns.svc.cluster.local
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

	// Create a stream.
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "mystream",
		Subjects: []string{"ORDERS.*"},
		Storage:  nats.FileStorage,
		MaxAge:   time.Hour * 1,
	})
	if err != nil {
		log.Fatalf("failed to create stream: %s\n", err)
	}

	// Create a consumer for the stream.
	_, err = js.AddConsumer("mystream", &nats.ConsumerConfig{
		Durable: "PULL_CONSUMER",
		AckPolicy: nats.AckExplicitPolicy,
	})
	if err != nil {
		log.Fatalf("failed to create consumer: %s\n", err)
	}
}
