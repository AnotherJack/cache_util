package cache_util

import (
	"context"
	"reflect"
	"time"
)

// 多级缓存，按传入的顺序排列优先级
type multiCacheStrategy struct {
	strategies []CacheStrategy
}

func (m *multiCacheStrategy) Put(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	var retErr error
	for _, s := range m.strategies {
		err := s.Put(ctx, key, value, expiration)
		if err != nil {
			retErr = err
		}
	}
	return retErr
}

func (m *multiCacheStrategy) Get(ctx context.Context, key string, resType reflect.Type) (interface{}, error) {
	for _, s := range m.strategies {
		v, err := s.Get(ctx, key, resType)
		if err == nil && !isNil(v) {
			return v, nil
		}
	}

	return nil, nil
}

func NewMultiCacheStrategy(strategies ...CacheStrategy) CacheStrategy {
	return &multiCacheStrategy{
		strategies: strategies,
	}
}
