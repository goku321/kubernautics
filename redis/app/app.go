package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	redis "github.com/go-redis/redis/v8"
)

func splitAndTrim(s, sep, toTrim string) []string {
	x := strings.Split(s, sep)
	for i := range x {
		x[i] = strings.Trim(x[i], toTrim)
	}
	return x
}

func writeToRedisList() error {
	addrString := os.Getenv("REDIS_ADDRESS")
	addrs := splitAndTrim(addrString, ",", " ")
	pass := os.Getenv("REDIS_PASSWORD")
	opts := redis.ClusterOptions{
		Addrs:    addrs,
		Password: pass,
	}
	client := redis.NewClusterClient(&opts)
	list := os.Getenv("LIST_NAME")
	itemCount, err := strconv.ParseInt(os.Getenv("NO_LIST_ITEMS_TO_WRITE"), 10, 32)
	if err != nil {
		return fmt.Errorf("number of items to write should be a number: %s", err.Error())
	}
	for i := 0; i < int(itemCount); i++ {
		x := client.LPush(context.Background(), list, i)
		if x.Err() != nil {
			return fmt.Errorf("failed to write to redis list: %s", x.Err().Error())
		}
	}
	return nil
}

func main() {
	// Print env and args.
	log.Println("REDIS_ADDRESS: ", os.Getenv("REDIS_ADDRESS"))
	log.Println("REDIS_PASSWORD: ", os.Getenv("REDIS_PASSWORD"))
	log.Println("LIST_NAME: ", os.Getenv("LIST_NAME"))
	log.Println("NO_LIST_ITEMS_TO_WRITE: ", os.Getenv("NO_LIST_ITEMS_TO_WRITE"))
	log.Println("Args: ", os.Args)

	action := ""
	if len(os.Args) > 0 {
		action = os.Args[1]
	}
	if action == "write" {
		err := writeToRedisList()
		if err != nil {
			log.Fatalf("write to redis list failed: %v\n", err)
		}
		log.Println("write to redis list is successful")
	} else if action == "read" {

	} else {
		log.Printf("unknown action: %s\n", action)
	}
}
