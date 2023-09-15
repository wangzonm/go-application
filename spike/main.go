package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"github.com/spf13/cast"
	"math/rand"
	"spike/cache"
	"spike/db"
	"time"
)

type Return struct {
	Code int
	Msg  string
	Data interface{}
}

var (
	COMMAND_NUM = 32
)

func main() {
	cache.Init()
	db.Init()
	engine := gin.Default()
	engine.GET("/v1/activity/init", ActivityInit)
	engine.GET("/v1/activity/sign_up", SignUpV3)
	engine.GET("/v1/activity/mutil_without_pipline", MutilOperationWithoutPipline)
	engine.GET("/v1/activity/mutil_with_pipline", MutilOperationWithPipline)
	engine.Run()
}

func ActivityInit(c *gin.Context) {
	info := cache.ActivityInfo{
		TopPkId:                    1,
		ShowStartTime:              0,
		ShowEndTime:                0,
		ShowStatus:                 0,
		ApplyStartTime:             0,
		ApplyEndTime:               0,
		PkSeconds:                  0,
		PunishSeconds:              0,
		AllowJoinCount:             0,
		CurrentJoinCount:           0,
		TopPkStatus:                0,
		CurrentRoundNumber:         0,
		CurrentRoundMemberCount:    0,
		CurrentRoundWinMemberCount: 0,
		NextRoundPkStartTime:       0,
		NextRoundReadyStartTime:    0,
		PkWinTime:                  0,
	}
	b, err := json.Marshal(&info)
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	var m map[string]interface{}
	if err = json.Unmarshal(b, &m); err != nil {
		c.AbortWithStatus(500)
		return
	}

	vals, err := cache.RedisClient.HMSet(cache.CACHE_ACTIVITY_INFO, m).Result()
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	c.JSON(200, fmt.Sprintf("init success, val=%+v", vals))
	return
}

//
// SignUpV1 scard和sadd是两个独立操作，非原子操作，实际元素会超过32个
// wrk -t10 -c1000 -d10s -s sign_up.lua http://127.0.0.1:8080
// @Summary SignUpV1
// @Tags
// @Produce json
// @Param c
// @Success 200 {object}
// @Failure 500
// @Router path [method]
func SignUpV1(c *gin.Context) {
	uid := c.Query("uid")
	//res := cache.RedisClient.Incr(cache.CACHE_ACTIVITY_SIGN_UP_COUNT)
	//count, err := res.Result()
	//if err != nil {
	//	c.AbortWithStatus(500)
	//	return
	//}
	//
	//if count >= 32 {
	//	c.AbortWithStatusJSON(401, "sign_up is full")
	//	return
	//}

	res := cache.RedisClient.SCard(cache.CACHE_ACTIVITY_SIGN_UP_SUCC)
	count, err := res.Result()
	if err != nil {
		c.AbortWithStatus(500)
		return
	}

	if count >= 32 {
		c.AbortWithStatusJSON(402, "sign_up is full")
		return
	}
	res = cache.RedisClient.SAdd(cache.CACHE_ACTIVITY_SIGN_UP_SUCC, uid)
	vals, err := res.Result()
	if err != nil {
		c.AbortWithStatus(500)
		return
	}
	c.JSON(200, fmt.Sprintf("sign_up success, val=%+v", vals))
	return
}

//
// SignUpV2 通过lua脚本保证scard和sadd操作的原子性，实际元素不会超过32

// wrk -t20 -c160 -d10s -s sign_up.lua http://127.0.0.1:8080
//Running 10s test @ http://127.0.0.1:8080
//20 threads and 160 connections
//Thread Stats   Avg      Stdev     Max   +/- Stdev
//Latency     5.76ms    2.72ms  52.03ms   81.13%
//Req/Sec     1.42k   277.55     2.41k    64.65%
//283265 requests in 10.02s, 22.97MB read
//Non-2xx or 3xx responses: 283233
//Requests/sec:  28283.25
//Transfer/sec:      2.29MB

//➜  spike git:(master) ✗ grep '| 401 |' log.log | wc -l
//282882
//➜  spike git:(master) ✗ grep '| 402 |' log.log | wc -l
//475
//➜  spike git:(master) ✗ grep dao_operation log.log| wc -l
//32

//wrk -t10 -c1000 -d10s -s sign_up.lua http://127.0.0.1:8080
//sign_up.lua: cannot open sign_up.lua: Too many open files
//Running 10s test @ http://127.0.0.1:8080
//10 threads and 1000 connections
//Thread Stats   Avg      Stdev     Max   +/- Stdev
//Latency     6.47ms   12.74ms 301.34ms   99.10%
//Req/Sec     3.66k     1.94k    8.48k    65.94%
//287195 requests in 10.06s, 23.28MB read
//Socket errors: connect 757, read 113, write 0, timeout 0
//Non-2xx or 3xx responses: 287163
//Requests/sec:  28534.04
//Transfer/sec:      2.31MB

// @Summary SignUpV2
// @Tags
// @Produce json
// @Param c
// @Success 200 {object}
// @Failure 500
// @Router path [method]
func SignUpV2(c *gin.Context) {
	uid := c.Query("uid")
	script := `
		local ex = redis.call('SCARD',KEYS[1])
        if ex >= 32
		then
			return -1
		else
			local res = redis.call('SADD',KEYS[1],ARGV[1]) 
			return res
		end
	`
	vals, err := cache.RedisClient.Eval(script, []string{cache.CACHE_ACTIVITY_SIGN_UP_SUCC}, uid).Result()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if vals.(int64) == -1 {
		c.AbortWithError(401, errors.New("sign_up is full"))
		return
	}

	var total int64
	total = 32
	if vals.(int64) == 1 {
		vals, err = cache.RedisClient.Incr(cache.CACHE_ACTIVITY_SIGN_UP_COUNT).Result()
		if err != nil {
			c.AbortWithError(500, err)
			return
		}

		fmt.Printf("dao_operation uid=%v vals=%v \n", uid, vals)
		affectRows, err := db.MySQLClient.Where("id = ? and count < ?", 1, 32).Incr("count", 1).Update(&db.Pk{})
		if err != nil {
			c.AbortWithError(501, err)
			return
		}

		fmt.Printf("sign_up_end_before uid=%v vals=%v affect_rows=%v \n", uid, vals, affectRows)
		if vals == total {
			fmt.Printf("sign_up_end uid=%v vals=%v \n", uid, vals)
		}
	} else {
		c.AbortWithError(402, errors.New("sign_up is full"))
		return
	}

	c.JSON(200, fmt.Sprintf("sign_up success, val=%+v", vals))
	return
}

//
// SignUpV3 有bug incr操作是原子的，最终不超过32个，但是导致sadd的数量小于32个
// @Summary SignUpV3
// @Tags
// @Produce json
// @Param c
// @Success 200 {object}
// @Failure 500
// @Router path [method]
func SignUpV3(c *gin.Context) {
	uid := c.Query("uid")
	ok, err := cache.RedisClient.SIsMember(cache.CACHE_ACTIVITY_SIGN_UP_SUCC, uid).Result()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if ok {
		c.AbortWithStatusJSON(402, "sign_up is full")
		return
	}

	vals, err := cache.RedisClient.Incr(cache.CACHE_ACTIVITY_SIGN_UP_COUNT).Result()
	if err != nil {
		c.AbortWithError(500, err)
		return
	}

	if vals > 32 {
		cache.RedisClient.Decr(cache.CACHE_ACTIVITY_SIGN_UP_COUNT)
	} else {
		if vals == 32 {
			c.AbortWithStatusJSON(401, "sign_up is full")
			return
		}

		fmt.Printf("sign_up_succ_uid=%v \n", uid)
		vals, err = cache.RedisClient.SAdd(cache.CACHE_ACTIVITY_SIGN_UP_SUCC, uid).Result()
		if err != nil {
			c.AbortWithStatus(500)
			return
		}

		c.JSON(200, fmt.Sprintf("sign_up success, val=%+v", vals))
		return
	}
}

func MutilOperationWithoutPipline(c *gin.Context) {
	testType := c.Query("type")
	start := time.Now()

	switch testType {
	case "list":
		RedisListWithOutPipline()
	case "hash":
		RedisHsetWithOutPipline()
	case "zset":
		RedisZsetWithOutPipline()
	}

	cost := time.Since(start).Microseconds()
	fmt.Printf("type=%v cost=%vus \n", testType, cost)
	c.JSON(200, "default")
	return
}

func MutilOperationWithPipline(c *gin.Context) {
	testType := c.Query("type")
	start := time.Now()

	switch testType {
	case "list":
		RedisListWithPipline()
	case "hash":
		RedisHsetWithPipline()
	case "zset":
		RedisZsetWithPipline()
	}

	cost := time.Since(start).Microseconds()
	fmt.Printf("type=%v cost=%vus \n", testType, cost)
	c.JSON(200, "default")
	return
}

func RedisListWithOutPipline() {
	for i := 0; i < COMMAND_NUM; i++ {
		_, err := cache.RedisClient.RPop(cache.CACHE_ACTIVITY_LIST).Result()
		if err != nil {
			fmt.Printf("rpop err=%v \n", err)
		}
	}
	return
}

func RedisListWithPipline() {
	pip := cache.RedisClient.Pipeline()
	for i := 0; i < COMMAND_NUM; i++ {
		pip.RPop(cache.CACHE_ACTIVITY_LIST)
	}
	pip.Exec()
	return
}

func RedisHsetWithOutPipline() {
	cache.RedisClient.HSet(cache.CACHE_ACTIVITY_HASH, "status", 1)
	cache.RedisClient.HSet(cache.CACHE_ACTIVITY_HASH, "current_number", 2)
	cache.RedisClient.HSet(cache.CACHE_ACTIVITY_HASH, "timestamp", time.Now().UnixMilli())
	cache.RedisClient.HSet(cache.CACHE_ACTIVITY_HASH, "round", 3)
	//cache.RedisClient.HMSet(cache.CACHE_ACTIVITY_HASH, map[string]interface{}{
	//	"status": 1,
	//	"current_number": 2,
	//	"timestamp": time.Now().UnixMilli(),
	//})
	return
}

func RedisHsetWithPipline() {
	pipe := cache.RedisClient.Pipeline()
	pipe.HSet(cache.CACHE_ACTIVITY_HASH, "status", 1)
	pipe.HSet(cache.CACHE_ACTIVITY_HASH, "current_number", 2)
	pipe.HSet(cache.CACHE_ACTIVITY_HASH, "timestamp", time.Now().UnixMilli())
	pipe.HSet(cache.CACHE_ACTIVITY_HASH, "round", 3)
	pipe.Exec()
	return
}

type Data struct {
	Uid       int    `json:"uid"`
	Name      string `json:"name"`
	Tag       int    `json:"tag"`
	Score     int    `json:"score"`
	Timestamp int64  `json:"timestamp"`
}

func RedisZsetWithOutPipline() {
	for i := 0; i < COMMAND_NUM; i++ {
		d := Data{
			Uid:       1000 + rand.Intn(1000000),
			Tag:       0,
			Score:     rand.Intn(100000),
			Timestamp: time.Now().UnixMilli(),
		}
		d.Name = fmt.Sprintf("nickname-%v", d.Uid)
		cache.RedisClient.ZAdd(cache.CACHE_ACTIVITY_ZSET, redis.Z{
			Score:  cast.ToFloat64(d.Score),
			Member: d.Uid,
		})
	}
}

func RedisZsetWithPipline() {
	pip := cache.RedisClient.Pipeline()
	for i := 0; i < COMMAND_NUM; i++ {
		d := Data{
			Uid:       1000 + rand.Intn(1000000),
			Tag:       0,
			Score:     rand.Intn(100000),
			Timestamp: time.Now().UnixMilli(),
		}
		d.Name = fmt.Sprintf("nickname-%v", d.Uid)
		pip.ZAdd(cache.CACHE_ACTIVITY_ZSET, redis.Z{
			Score:  cast.ToFloat64(d.Score),
			Member: d.Uid,
		})
	}
	pip.Exec()
}
