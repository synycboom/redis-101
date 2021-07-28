package simplequeue

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis/v8"
)

type Queue struct {
	name string
	rdb  *redis.ClusterClient
}

func New(clusterAddr []string, password, name string) *Queue {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    clusterAddr,
		Password: password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(errors.Wrap(err, "[New]: unable to ping"))
	}

	return &Queue{
		rdb:  rdb,
		name: name,
	}
}

func (q *Queue) Produce(item string) int64 {
	ctx := context.Background()

	res, err := q.rdb.LPush(ctx, q.name, item).Result()
	if err != nil {
		panic(err)
	}

	return res
}

func (q *Queue) Consume(timeout time.Duration) []string {
	ctx := context.Background()

	res, err := q.rdb.BRPop(ctx, timeout, q.name).Result()
	if err != nil {
		panic(err)
	}

	return res
}
