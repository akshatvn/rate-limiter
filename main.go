package main

import (
	"fmt"
	"github.com/akshatvn/rate-limiter/ratelimiter"
	"gopkg.in/redis.v5"
	"math/rand"
	"time"
)
// short program demonstrating how to use the rate limiter package
func main() {
	// initialize a redis client
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	// define limits per user/tenant/third-party
	limitMap := map[string]uint64{
		"Raamu":  20,
		"Shaamu": 40,
	}

	// initialize the ratelimiter
	r := ratelimiter.NewRedisRateLimiter(5, limitMap, redisClient)

	startTime := time.Now().UnixNano()
	rand.Seed(startTime)

	// And do some random testing and printing
	for i := 0; i < 100; i++ {
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(500)))
		nowTime := time.Now().UnixNano()

		// result is false if limit is not breached.
		if result, _ := r.IsLimitBreached("Raamu"); result {
			fmt.Printf("\t%d\tlimit breached!\ttime passed=%dms \n", i, (nowTime-startTime)/1000000)
		} else {
			fmt.Printf("\t%d\tDo your thing\ttime passed=%dms \n", i, (nowTime-startTime)/1000000)
		}

	}
}

