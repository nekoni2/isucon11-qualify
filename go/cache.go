package main

import (
	"context"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var redisClient *redis.Client

// redis cache key
func cacheKey(arr ...string) string {
	return strings.Join(arr, ":")
}

// create redis client
func initRedis() {
	addr := getEnv("REDIS_ADDR", "192.168.0.12:6379")
	if addr == "" {
		logger.Error("Redis addr env empty")
		return
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := rdb.Ping(ctx).Result(); err != nil {
		logger.Errorf("Redis ping err: %s", err)
	}

	redisClient = rdb
}

