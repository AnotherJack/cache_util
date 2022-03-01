package cache_util

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/go-redis/redis/v8"
)

func TestCacheUtil_GetWithCache(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	myRC := &myRedisClient{rdb: rdb}
	cUtil := NewCacheUtil(NewRedisCacheStrategy(myRC))
	for i := 0; i < 10; i++ {
		var numInt int64
		num, err := cUtil.GetWithCache(context.Background(), "asd", reflect.TypeOf(numInt), time.Second*200, func() (interface{}, error) {
			r := RandInt64WithinRange(1, 100)
			fmt.Printf("f executed, r is %v \n", r)
			// 生成一个随机数
			return r, nil
		})
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("loop i:%v, num:%v \n", i, num)
		}
	}

}

func TestCacheQuerier_GetWithCache(t *testing.T) {
	cUtil := NewCacheUtil(NewMapCacheStrategy())
	q := NewCahceQuerier[int64](cUtil)
	for i := 0; i < 10; i++ {
		// var numInt int64
		numInt, err := q.GetWithCache(context.Background(), "asd", time.Second, func() (int64, error) {
			r := RandInt64WithinRange(1, 100)
			fmt.Printf("f executed, r is %v \n", r)
			// 生成一个随机数
			return r, nil
		})
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("loop i:%v, num:%v \n", i, numInt)
		}
	}
}

func TestCacheQuerier_GetWithRedisCache(t *testing.T) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	myRC := &myRedisClient{rdb: rdb}
	cUtil := NewCacheUtil(NewRedisCacheStrategy(myRC))
	q := NewCahceQuerier[*User](cUtil)
	for i := 0; i < 10; i++ {
		// var numInt int64
		uId := fmt.Sprint(i%3 + 1)
		key := fmt.Sprintf("key:user:%v", uId)
		user, err := q.GetWithCache(context.Background(), key, time.Second*100, func() (*User, error) {
			fmt.Printf("read from rpc uId:%v \n", uId)
			return getUserInfo(uId)
		})
		if err != nil {
			panic(err)
		} else {
			fmt.Printf("loop i:%v, uid:%v, uName:%v \n", i, uId, user.Name)
		}
	}
}

type User struct {
	Id   string `json:"id"`
	Name string `json:"name"`
}

func getUserInfo(userId string) (*User, error) {
	u := &User{Id: userId}
	switch userId {
	case "1":
		u.Name = "Jack"
	case "2":
		u.Name = "Tom"
	case "3":
		u.Name = "Bob"
	default:
		u.Name = "Jack"
	}
	return u, nil
}

type myRedisClient struct {
	rdb *redis.Client
}

func (m *myRedisClient) Get(ctx context.Context, key string) (string, error) {
	cmd := m.rdb.Get(ctx, key)
	return cmd.Result()
}

func (m *myRedisClient) Set(ctx context.Context, key string, value string, expiration time.Duration) error {
	cmd := m.rdb.Set(ctx, key, value, expiration)
	return cmd.Err()
}
