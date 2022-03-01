package cache_util

import (
	"context"
	"reflect"
	"time"
)

type CacheUtil interface {
	GetWithCache(ctx context.Context, key string, resType reflect.Type, expiration time.Duration, f func() (interface{}, error)) (interface{}, error)
}

type CacheStrategy interface {
	Put(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, resType reflect.Type) (interface{}, error)
}

type cacheUtil struct {
	cacheStrategy CacheStrategy
}

func (c *cacheUtil) GetWithCache(ctx context.Context, key string, resType reflect.Type, expiration time.Duration, f func() (interface{}, error)) (interface{}, error) {
	// 先从缓存拿
	cachedRes, err := c.cacheStrategy.Get(ctx, key, resType)
	if err == nil && !isNil(cachedRes) {
		return cachedRes, nil
	}

	// 从f中拿
	res, err := f()
	if err != nil {
		return nil, err
	}
	if isNil(res) {
		return res, nil
	}

	// 存缓存
	c.cacheStrategy.Put(ctx, key, res, expiration)

	return res, nil
}

func NewCacheUtil(strategy CacheStrategy) CacheUtil {
	return &cacheUtil{
		cacheStrategy: strategy,
	}
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
