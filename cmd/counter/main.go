package main

import (
	"log"
	"os"
	"strconv"
	"strings"
	"sync"

	"github.com/synycboom/redis-101/transaction"
)

var (
	isAtomicCounter bool
	counterName     string
	maxRetries      int
	redisPassword   string
	redisAddress    []string
)

func init() {
	isAtomicCounter, _ = strconv.ParseBool(os.Getenv("IS_ATOMIC_COUNTER"))
	maxRetries, _ = strconv.Atoi(os.Getenv("COUNTER_MAX_RETRY"))
	counterName = os.Getenv("COUNTER_NAME")
	if counterName == "" {
		panic("[init]: counterName is required")
	}

	redisPassword = os.Getenv("REDIS_PASSWORD")
	if redisPassword == "" {
		panic("[init]: REDIS_PASSWORD is required")
	}

	redisAddress = strings.Split(os.Getenv("REDIS_ADDRESS"), ",")
	if len(redisAddress) == 0 {
		panic("[init]: CLUSTER_ADDR is required")
	}
}

func main() {
	counter := transaction.NewCounter(redisAddress, redisPassword)
	initialVal := counter.Value(counterName)

	log.Printf("[main]: initial value of counter '%s' is %d\n", counterName, initialVal)

	var wg sync.WaitGroup
	for i := 1; i <= 20; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			if isAtomicCounter {
				counter.AtomicIncrease(counterName, maxRetries)
			} else {
				counter.Increase(counterName)
			}
		}()
	}

	wg.Wait()

	latestValue := counter.Value(counterName)

	log.Printf("[main]: lastest value of counter '%s' is %d\n", counterName, latestValue)
	log.Printf("[main]: the diff of the initial and latest of counter '%s' is %d\n", counterName, latestValue-initialVal)
}
