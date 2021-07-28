package pubsub

import (
	"context"
	"log"
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis/v8"
)

type PubSub struct {
	rdb *redis.ClusterClient
}

func New(clusterAddr []string, password string) *PubSub {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    clusterAddr,
		Password: password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(errors.Wrap(err, "[New]: unable to ping"))
	}

	return &PubSub{
		rdb: rdb,
	}
}

func (p *PubSub) Subscribe(channels ...string) {
	ctx := context.Background()

	pubsub := p.rdb.Subscribe(ctx, channels...)

	if _, err := pubsub.Receive(ctx); err != nil {
		panic(errors.Wrap(err, "[PubSub.Subscribe]: unable to subscribe"))
	}

	ch := pubsub.Channel()

	for msg := range ch {
		log.Printf("[Pubsub.Subscribe]: message '%s' from channel '%s'", msg.Payload, msg.Channel)
	}
}

func (p *PubSub) Publish(channel, message string) {
	ctx := context.Background()

	if err := p.rdb.Publish(ctx, channel, message).Err(); err != nil {
		panic(errors.Wrap(err, "[PubSub.Publish]: unable to publish"))
	}
}
