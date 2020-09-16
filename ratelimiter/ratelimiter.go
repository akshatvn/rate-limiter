package ratelimiter

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path"
	"runtime"
	"time"

	"gopkg.in/redis.v5"
)

// Api so you can extend to other implementations
type Api interface {
	IsLimitBreached() (bool, error)
}

type RedisRateLimiter struct {
	windowSize  int               // seconds
	limitMap    map[string]uint64 // maintain map of key to rate limit
	redisClient *redis.Client     // the active redis connection which stores all timestamps of requests that went through
}

// NewRedisRateLimiter returns a new rate limiter with the provided config
func NewRedisRateLimiter(windowSize int, limitMap map[string]uint64, redisClient *redis.Client) *RedisRateLimiter {
	if windowSize == 0 {
		windowSize = 60 // seconds
	}
	return &RedisRateLimiter{
		windowSize:  windowSize,
		limitMap:    limitMap,
		redisClient: redisClient,
	}
}

// IsLimitBreached is the core function that determines whether a limit is breached for a given key
func (r *RedisRateLimiter) IsLimitBreached(key string) (bool, error) {
	limit := r.limitMap[key]
	windowSize := r.windowSize // seconds
	filename := getFilePath()
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		_ = fmt.Errorf("could not read file %s: %+v", filename, err)
		return false, errors.New("error reading the lua file")
	}
	luaScript := string(content)
	cmd := r.redisClient.Eval(luaScript, []string{key}, time.Now().Unix(), windowSize, limit)
	vals, err := cmd.Result()
	fmt.Printf("%+v left:", vals)
	if valInt, ok := vals.(int64); ok {
		return valInt == 0, nil
	}
	return false, errors.New("the answer was not an integer")
}

// getFilePath is a helper function to get the path of supporting lua script
func getFilePath() string {
	_, filename, _, _ := runtime.Caller(1)
	return path.Join(path.Dir(filename), "/ratelimiter.lua")
}
