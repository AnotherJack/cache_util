package cache_util

import (
	"math/rand"
	"time"
)

func RandInt64WithinRange(from int64, to int64) int64 {
	rand.Seed(time.Now().UnixNano())
	randNum := rand.Int63n(to-from+1) + from
	return randNum
}
