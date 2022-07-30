package main

import (
	"log"
	"strconv"

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

	for i := 1; i <= 1000; i++ {
		_, err = js.Publish("ORDERS.scratch", []byte("order "+strconv.Itoa(i)))
		if err == nil {
			log.Printf("published message: %d\n", i)
		}
	}

	nc.Close()
}
