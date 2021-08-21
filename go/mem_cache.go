package main

import (
	"context"
	"errors"
	"strings"

	"github.com/dgraph-io/ristretto"
	jsoniter "github.com/json-iterator/go"
)

var cache *ristretto.Cache

func init() {
	var err error
	cache, err = ristretto.NewCache(&ristretto.Config{
		NumCounters: 1e7,     // Num keys to track frequency of (10M).
		MaxCost:     1 << 30, // Maximum cost of cache (1GB).
		BufferItems: 64,      // Number of keys per Get buffer.
	})
	if err != nil {
		panic(err)
	}
}

func getCacheFallback(key interface{}, fallback func() (interface{}, error)) (interface{}, error) {
	data, ok := cache.Get(key)
	if ok {
		return data, nil
	}
	var err error
	data, err = fallback()
	if err != nil {
		return nil, err
	}
	cache.Set(key, data, 1)
	return data, nil
}

func cacheGetIsuCondition(id string) (*IsuCondition, error) {
	id = strings.TrimPrefix(id, cachePrefixIsuContidion)
	data, err := getCacheFallback(cachePrefixIsuContidion + id, func() (interface{}, error) {
		jsonStr, err := redisClient.Get(context.Background(), cachePrefixIsuContidion + id).Result()
		if err != nil {
			return nil, err
		}
		obj := new(IsuCondition)
		if err := jsoniter.UnmarshalFromString(jsonStr, &obj); err != nil {
			return nil, err
		}
		return obj, nil
	})
	if err != nil {
		return nil, err
	}
	if val, ok := data.(*IsuCondition); ok {
		return val, nil
	}
	return nil, errors.New("val parse err")
}
