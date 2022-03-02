package cache_util

import (
	"context"
	"encoding/json"
	"reflect"
	"time"
)

type redisCacheStrategy struct {
	redisClient RedisClient
}

func (r *redisCacheStrategy) Put(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonStr, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return r.redisClient.Set(ctx, key, string(jsonStr), expiration)
}

func (r *redisCacheStrategy) Get(ctx context.Context, key string, resType reflect.Type) (interface{}, error) {
	cachedJson, err := r.redisClient.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	resPtrValue := reflect.New(resType)

	err = json.Unmarshal([]byte(cachedJson), resPtrValue.Interface())
	if err != nil {
		return nil, err
	}
	return resPtrValue.Elem().Interface(), nil
}

type RedisClient interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

func NewRedisCacheStrategy(redisClient RedisClient) CacheStrategy {
	return &redisCacheStrategy{
		redisClient: redisClient,
	}
}
