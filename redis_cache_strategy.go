package cache_util

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

type redisCacheStrategy struct {
	redisClient RedisClient
}

func (r *redisCacheStrategy) Put(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	encodedStr, err := r.encode(ctx, value)
	if err != nil {
		return fmt.Errorf("encode error:%+v", err)
	}
	return r.redisClient.Set(ctx, key, encodedStr, expiration)
}

func (r *redisCacheStrategy) Get(ctx context.Context, key string, resType reflect.Type) (interface{}, error) {
	rawStr, err := r.redisClient.Get(ctx, key)
	if err != nil {
		return nil, err
	}

	res, err := r.decode(ctx, rawStr, resType)
	if err != nil {
		return nil, fmt.Errorf("decode error:%+v", err)
	}

	return res, nil
}

func (r *redisCacheStrategy) encode(ctx context.Context, value interface{}) (string, error) {
	switch s := value.(type) {
	case int:
		return strconv.Itoa(s), nil
	case int8:
		return strconv.FormatInt(int64(s), 10), nil
	case int16:
		return strconv.FormatInt(int64(s), 10), nil
	case int32:
		return strconv.Itoa(int(s)), nil
	case int64:
		return strconv.FormatInt(s, 10), nil
	case float32:
		return strconv.FormatFloat(float64(s), 'f', -1, 32), nil
	case float64:
		return strconv.FormatFloat(s, 'f', -1, 64), nil
	case string:
		return s, nil
	case bool:
		return strconv.FormatBool(s), nil
	default:
		// encode with json
		jsonStr, err := json.Marshal(value)
		if err != nil {
			return "", err
		}
		return string(jsonStr), nil
	}
}

func (r *redisCacheStrategy) decode(ctx context.Context, rawStr string, t reflect.Type) (interface{}, error) {
	switch t.Kind() {
	case reflect.Int:
		return strconv.Atoi(rawStr)
	case reflect.Int8:
		i, err := strconv.ParseInt(rawStr, 10, 8)
		if err != nil {
			return int8(0), err
		}
		return int8(i), nil
	case reflect.Int16:
		i, err := strconv.ParseInt(rawStr, 10, 16)
		if err != nil {
			return int16(0), err
		}
		return int16(i), nil
	case reflect.Int32:
		i, err := strconv.ParseInt(rawStr, 10, 32)
		if err != nil {
			return int32(0), err
		}
		return int32(i), nil
	case reflect.Int64:
		i, err := strconv.ParseInt(rawStr, 10, 64)
		if err != nil {
			return int64(0), err
		}
		return i, nil
	case reflect.Float32:
		f, err := strconv.ParseFloat(rawStr, 32)
		if err != nil {
			return float32(0), err
		}
		return float32(f), nil
	case reflect.Float64:
		f, err := strconv.ParseFloat(rawStr, 64)
		if err != nil {
			return float64(0), err
		}
		return f, nil
	case reflect.String:
		return rawStr, nil
	case reflect.Bool:
		return strconv.ParseBool(rawStr)
	default:
		// decode with json
		resPtrValue := reflect.New(t)
		err := json.Unmarshal([]byte(rawStr), resPtrValue.Interface())
		if err != nil {
			return nil, err
		}
		return resPtrValue.Elem().Interface(), nil
	}
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
