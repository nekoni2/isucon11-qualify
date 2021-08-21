package main

import (
	"context"
	"runtime/debug"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
	jsoniter "github.com/json-iterator/go"
	"github.com/spf13/cast"
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

	redisClient = rdb
}

// cache a DB Row and encode in JSON
func cacheRow(prefix string, id interface{}, row interface{}, r redis.Pipeliner) {
	val, _ := jsoniter.MarshalToString(row)
	var err error
	if r != nil {
		err = r.Set(context.Background(), prefix+cast.ToString(id), val, 0).Err()
	} else {
		err = redisClient.Set(context.Background(), prefix+cast.ToString(id), val,
			0).Err()
	}
	if err != nil {
		logger.Errorf("redis cache row err: %s", err)
	}
}

// fetch one from cache and decode JSON
func fetchCacheRow(prefix string, id interface{}, data interface{}) error {
	val, err := redisClient.Get(context.Background(), prefix+cast.ToString(id)).Result()
	if err != nil {
		return err
	}
	if err := jsoniter.UnmarshalFromString(val, data); err != nil {
		return err
	}
	return nil
}

var (
	cachePrefixIsuContidion = "isucondition:id:"
)

func tablesCache() {
    logger.Info("table cache start")
	defer func() {
		if err := recover(); err != nil {
			logger.With("stack", string(debug.Stack())).Errorf("table cache err: %s", err)
		}
	}()
	
	timeStart := time.Now()
	// cache `isu_condition` table
	{
		query := "SELECT * FROM isu_condition"
		data := make([]*IsuCondition, 0, 3200)
		if err := db.Select(&data, query); err != nil {
			logger.Errorf("isu_condition table query err: %s", err)
		}
        logger.Infof("total isu_condition items: %v, loaded", len(data))

		pipe := redisClient.Pipeline()
		for _, row := range data {
			cacheRow(cachePrefixIsuContidion, row.ID, row, pipe)
		}
		if err := cacheIsuCondition(pipe, data...); err != nil {
			logger.Errorf("cache chair err: %s", err)
		}
		if _, err := pipe.Exec(context.Background()); err != nil {
			logger.Errorf("redis cache estate err: %s", err)
		} else {
            logger.Infof("redis cache estate success !!!")
        }
	}

	timeDuration := time.Since(timeStart)
	logger.Infof("redis cache exec time duration: %s", timeDuration.String())
}

// jia-uuid < 1 --- n > db id
func cacheIsuCondition(pipe redis.Pipeliner, arr ...*IsuCondition) error {
	ctx := context.Background()
	for _, row := range arr {
		pipe.ZAdd(ctx, cacheKey("isucondition", "timestamp", row.JIAIsuUUID), &redis.Z{
			Score: float64(row.Timestamp.Unix()),
			Member: row.ID, // db key
		})
	}
	return nil
}
