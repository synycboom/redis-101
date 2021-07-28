package main

import (
	"os"
	"strings"

	pubsub "github.com/synycboom/redis-101/pub-sub"
)

var (
	mode          string
	message       string
	channels      []string
	redisPassword string
	redisAddress  []string
)

func init() {
	mode = os.Getenv("MODE")
	if !(mode == "publisher" || mode == "subscriber") {
		panic("MODE should be 'publisher' or 'subscriber'")
	}

	channels = strings.Split(os.Getenv("CHANNELS"), ",")
	if len(channels) == 0 {
		panic("CHANNELS is required")
	}

	message = os.Getenv("MESSAGE")
	if mode == "publisher" && message == "" {
		panic("MESSAGE is required when using publisher mode")
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
	if mode == "publisher" {
		publisher := pubsub.New(redisAddress, redisPassword)
		for _, ch := range channels {
			publisher.Publish(ch, message)
		}
	} else {
		subscriber := pubsub.New(redisAddress, redisPassword)
		subscriber.Subscribe(channels...)
	}
}
