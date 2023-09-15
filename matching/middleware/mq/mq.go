package mq

import (
	"github.com/go-redis/redis"
	"matching/middleware"
)

// SendCancelResult 撤单消息
func SendCancelResult(symbol, orderId string, ok bool) {
	values := map[string]interface{}{"orderId": orderId, "ok": ok}
	a := &redis.XAddArgs{
		Stream:       "matching:cancelresults:" + symbol,
		MaxLenApprox: 1000,
		Values:       values,
	}
	middleware.RedisClient.XAdd(a)
}

// SendTrade 成交消息
func SendTrade(symbol string, trade map[string]interface{}) {
	a := &redis.XAddArgs{
		Stream:       "matching:trades:" + symbol,
		MaxLenApprox: 1000,
		Values:       trade,
	}
	middleware.RedisClient.XAdd(a)
}
