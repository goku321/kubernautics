package main

import (
	"log"
	"os"
	"strconv"

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

	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("failed to create JS context: %s", err)
	}

	messagesCount, err := strconv.ParseInt(os.Getenv("NUM_MESSAGES"), 10, 32)
	if err != nil {
		log.Fatalf("number of messages to write should be a number: %s", err.Error())
	}
	for i := 1; i <= int(messagesCount); i++ {
		_, err = js.Publish("ORDERS.scratch", []byte("order "+strconv.Itoa(i)))
		if err == nil {
			log.Printf("published message: %d\n", i)
		}
	}

	nc.Close()
}
