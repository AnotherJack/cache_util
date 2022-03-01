package cache_util

import (
	"context"
	"errors"
	"reflect"
	"time"
)

type CacheQuerier[T any] struct {
	cUtil CacheUtil
}

func NewCahceQuerier[T any](cUtil CacheUtil) *CacheQuerier[T] {
	return &CacheQuerier[T]{
		cUtil: cUtil,
	}
}

func (q *CacheQuerier[T]) GetWithCache(ctx context.Context, key string, expiration time.Duration, f func() (T, error)) (T, error) {
	var t T
	v, err := q.cUtil.GetWithCache(ctx, key, reflect.TypeOf(t), expiration, func() (interface{}, error) {
		return f()
	})
	var zero T
	if err != nil {
		return zero, err
	}
	if isNil(v) {
		return zero, nil
	}
	tv, ok := v.(T)
	if ok {
		return tv, nil
	}

	// 按指针强转一下
	tv2, ok := v.(*T)
	if ok {
		return *tv2, nil
	}

	return zero, errors.New("type cast error")
}
