package cache_util

import (
	"context"
	"reflect"
	"sync"
	"time"
)

type mapCacheStrategy struct {
	cacheMap *sync.Map
}

func (m *mapCacheStrategy) Put(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	m.cacheMap.Store(key, value)
	return nil
}

func (m *mapCacheStrategy) Get(ctx context.Context, key string, resType reflect.Type) (interface{}, error) {
	v, ok := m.cacheMap.Load(key)
	if ok {
		return v, nil
	} else {
		return nil, nil
	}
}

func NewMapCacheStrategy() CacheStrategy {
	return &mapCacheStrategy{
		cacheMap: &sync.Map{},
	}
}
