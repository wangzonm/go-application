package main

import (
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis"
	"github.com/spf13/cast"
	"math/rand"
	"spike/cache"
	"testing"
	"time"
)

func init() {
	cache.Init()
}

func TestBuildCacheList(t *testing.T) {
	type data struct {
		Uid       int    `json:"uid"`
		Name      string `json:"name"`
		Tag       int    `json:"tag"`
		Score     int    `json:"score"`
		Timestamp int64  `json:"timestamp"`
	}
	for i := 1; i <= 1000000; i++ {
		d := data{
			Uid:       1000 + i,
			Name:      fmt.Sprintf("nickname-%v", i),
			Tag:       8,
			Score:     rand.Intn(1000000),
			Timestamp: time.Now().UnixMilli(),
		}
		b, _ := json.Marshal(&d)
		cache.RedisClient.LPush(cache.CACHE_ACTIVITY_LIST, string(b))
	}
}

func TestBuildCacheZset(t *testing.T) {
	type data struct {
		Uid       int    `json:"uid"`
		Name      string `json:"name"`
		Tag       int    `json:"tag"`
		Score     int    `json:"score"`
		Timestamp int64  `json:"timestamp"`
	}
	for i := 1; i <= 1000000; i++ {
		d := data{
			Uid:       1000 + i,
			Name:      fmt.Sprintf("nickname-%v", i),
			Tag:       8,
			Score:     rand.Intn(1000000),
			Timestamp: time.Now().UnixMilli(),
		}
		cache.RedisClient.ZAdd(cache.CACHE_ACTIVITY_ZSET, redis.Z{
			Score:  cast.ToFloat64(d.Score),
			Member: d.Uid,
		})
	}
}
