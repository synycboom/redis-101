package transaction

import (
	"context"
	"log"
	"math/rand"
	"time"

	"github.com/pkg/errors"

	"github.com/go-redis/redis/v8"
)

type Counter struct {
	rdb *redis.ClusterClient
}

func NewCounter(clusterAddr []string, password string) *Counter {
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    clusterAddr,
		Password: password,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(10)*time.Second)
	defer cancel()

	if _, err := rdb.Ping(ctx).Result(); err != nil {
		panic(errors.Wrap(err, "[NewCounter]: unable to ping"))
	}

	return &Counter{
		rdb: rdb,
	}
}

func (a *Counter) Value(counterName string) int {
	ctx := context.Background()
	res, err := a.rdb.Get(ctx, counterName).Int()
	if err == redis.Nil {
		return 0
	}
	if err != nil {
		panic(errors.Wrap(err, "[Counter.Value]: unable to get a counter"))
	}

	return res
}

func (a *Counter) Increase(counterName string) {
	ctx := context.Background()

	n := a.Value(counterName)
	if err := a.rdb.Set(ctx, counterName, n+1, 0).Err(); err != nil {
		panic(errors.Wrap(err, "[Counter.Increase]: unable to set a counter"))
	}
}

func (c *Counter) AtomicIncrease(counterName string, maxRetries int) {
	ctx := context.Background()
	txf := func(tx *redis.Tx) error {
		n, err := tx.Get(ctx, counterName).Int()
		if err != nil && err != redis.Nil {
			return err
		}

		n++

		log.Println(n)

		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, counterName, n, 0)

			return nil
		})

		return err
	}

	for retries := 0; retries < maxRetries; retries++ {
		err := c.rdb.Watch(ctx, txf, counterName)
		if err == nil {
			return
		}
		if err == redis.TxFailedErr {
			log.Println("[Counter.AtomicIncrease]: optimistic lock failure, sleep a bit and keep retrying...")

			randomSleep()

			continue
		}

		panic(errors.Wrap(err, "[Counter.AtomicIncrease]: unknown error occurs"))
	}

	panic("[Counter.AtomicIncrease]: increment reached maximum number of retries")
}

// ransomSleep sleep randomly for 0-10 ms
func randomSleep() {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10)
	time.Sleep(time.Duration(n) * time.Millisecond)
}
