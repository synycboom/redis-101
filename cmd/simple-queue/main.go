package main

import (
	"log"
	"os"
	"strings"
	"time"

	simplequeue "github.com/synycboom/redis-101/simple-queue"
)

var (
	mode          string
	message       string
	queue         string
	redisPassword string
	redisAddress  []string
)

func init() {
	mode = os.Getenv("MODE")
	if !(mode == "producer" || mode == "consumer") {
		panic("MODE should be 'producer' or 'consumer'")
	}

	queue = os.Getenv("QUEUE")
	if queue == "" {
		panic("QUEUE is required")
	}

	message = os.Getenv("MESSAGE")
	if mode == "producer" && message == "" {
		panic("MESSAGE is required when using producer mode")
	}

	redisPassword = os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		panic("REDIS_PASSWORD is required")
	}

	redisAddress = strings.Split(os.Getenv("REDIS_ADDRESS"), ",")
	if len(redisAddress) == 0 {
		panic("CLUSTER_ADDR is required")
	}
}

func main() {
	if mode == "producer" {
		producer := simplequeue.New(redisAddress, redisPassword, queue)
		producer.Produce(message)
	} else {
		consumer := simplequeue.New(redisAddress, redisPassword, queue)

		for {
			results := consumer.Consume(time.Duration(0))

			log.Printf("[main]: message '%s' from queue '%s'\n", results[1], results[0])
		}
	}
}
