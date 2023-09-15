package cache

import (
	"github.com/go-redis/redis"
	"log"
)

var RedisClient *redis.Client

const (
	CACHE_ACTIVITY_INFO          = "activity:i"        //活动详情
	CACHE_ACTIVITY_SIGN_UP_COUNT = "activity:s:u:incr" //报名人数统计
	CACHE_ACTIVITY_SIGN_UP_SUCC  = "activity:s:u:succ" //活动报名成功名单set
	CACHE_ACTIVITY_HASH          = "activity:hash"
	CACHE_ACTIVITY_LIST          = "activity:list"
	CACHE_ACTIVITY_ZSET          = "activity:zset"
)

//
//  ActivityInfo
//  @Description: 活动详情
//
type ActivityInfo struct {
	TopPkId                    int64 `json:"top_pk_id"`
	ShowStartTime              int   `json:"show_start_time"`
	ShowEndTime                int   `json:"show_end_time"`
	ShowStatus                 int8  `json:"show_status"` // 1显示中 0结束显示
	ApplyStartTime             int   `json:"apply_start_time"`
	ApplyEndTime               int   `json:"apply_end_time"`
	PkSeconds                  int   `json:"pk_seconds"`
	PunishSeconds              int   `json:"punish_seconds"`
	AllowJoinCount             int   `json:"allow_join_count"`
	CurrentJoinCount           int   `json:"current_join_count"`
	TopPkStatus                int8  `json:"top_pk_status"` // 0未开始报名 1报名中 2准备pk中 3pk中 4结束
	CurrentRoundNumber         int   `json:"current_round_number"`
	CurrentRoundMemberCount    int   `json:"current_round_member_count"`
	CurrentRoundWinMemberCount int   `json:"current_round_win_member_count"`
	NextRoundPkStartTime       int   `json:"next_round_pk_start_time"`
	NextRoundReadyStartTime    int   `json:"next_round_ready_start_time"`
	PkWinTime                  int   `json:"pk_win_time"`
}

func Init() {
	addr := "127.0.0.1:6379"
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err := RedisClient.Ping().Result()
	if err != nil {
		panic(err)
	} else {
		log.Printf("Connected to redis: %s", addr)
	}
}
