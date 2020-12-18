package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/go-redis/redis"
)

func splitAndTrim(s, sep, toTrim string) []string {
	x := strings.Split(s, sep)
	for i := range x {
		x[i] = strings.Trim(x[i], toTrim)
	}
	return x
}

func parseAddress() []string {
	addrString := os.Getenv("REDIS_ADDRESSES")
	if len(addrString) != 0 {
		return splitAndTrim(addrString, ",", " ")
	}
	hostString := os.Getenv("REDIS_HOSTS")
	portString := os.Getenv("REDIS_PORTS")
	hosts := splitAndTrim(hostString, ",", " ")
	ports := splitAndTrim(portString, ",", " ")
	addrs := []string{}
	if len(hosts) != len(ports) {
		return addrs
	}
	for i := range hosts {
		addrs = append(addrs, fmt.Sprintf("%s:%s", hosts[i], ports[i]))
	}
	return addrs
}

func redisStreamConsumer() error {
	addrs := parseAddress()
	pass := os.Getenv("REDIS_PASSWORD")
	opts := redis.ClusterOptions{
		Addrs:    addrs,
		Password: pass,
	}
	client := redis.NewClusterClient(&opts)
	stream := os.Getenv("REDIS_STREAM_NAME")
	_ = os.Getenv("REDIS_STREAM_CONSUMER_GROUP_NAME")

	for {
		len, err := client.XLen(context.Background(), stream).Result()
		if err != nil {
			return err
		}
		if len > 0 {
			x := client.XRead(context.Background(), &redis.XReadArgs{
				Streams: []string{stream},
				Count:   1,
				Block:   200 * time.Millisecond,
			})
			if x.Err() != nil {
				return fmt.Errorf("failed to read from redis stream: %s", x.Err().Error())
			}
		}
	}
	return nil
}

func redisStreamProducer() error {
	return nil
}

func main() {
	// Print env and args.
	log.Println("REDIS_ADDRESSES: ", parseAddress())
	log.Println("REDIS_PASSWORD: ", os.Getenv("REDIS_PASSWORD"))
	log.Println("REDIS_STREAM_NAME: ", os.Getenv("REDIS_STREAM_NAME"))
	log.Println("REDIS_STREAM_CONSUMER_GROUP_NAME: ", os.Getenv("REDIS_STREAM_CONSUMER_GROUP_NAME"))
	log.Println("NUM_MESSAGES: ", os.Getenv("NUM_MESSAGES"))
	log.Println("Args: ", os.Args)

	mode := ""
	if len(os.Args) > 0 {
		mode = os.Args[1]
	}
	if mode == "consumer" {
		if err := redisStreamConsumer(); err != nil {
			log.Fatalf("read from redis stream failed: %v\n", err)
		}
		log.Println("read from redis stream is successful")
	} else if mode == "producer" {
		if err := redisStreamProducer(); err != nil {
			log.Fatalf("write to redis stream failed: %v\n", err)
		}
		log.Println("write to redis stream is successful")
	} else {
		log.Printf("unknown mode: %s\n", mode)
	}
}
