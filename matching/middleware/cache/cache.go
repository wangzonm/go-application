package cache

import (
	"github.com/go-redis/redis"
	"github.com/shopspring/decimal"
	"matching/middleware"
	"strconv"
	"strings"
)

// Redis缓存设计目的
//1.请求去重，避免重复提交相同订单；
//2.恢复数据，即程序重启后能恢复所有数据。

func GetOrder(symbol, orderId string) map[string]string {
	orderData := strings.Split(orderId, ":")
	key := "matching:order:" + symbol + ":" + orderData[0] + ":" + orderData[1]
	return middleware.RedisClient.HGetAll(key).Val()
}

// SaveSymbol 每个symbol设置唯一缓存
func SaveSymbol(symbol string) {
	key := "matching:symbols"
	middleware.RedisClient.SAdd(key, symbol)
}

func RemoveSymbol(symbol string) {
	key := "matching:symbols"
	middleware.RedisClient.SRem(key, symbol)
}

func GetSymbols() []string {
	key := "matching:symbols"
	return middleware.RedisClient.SMembers(key).Val()
}

func OrderExist(symbol, orderId, action string) bool {
	key := "matching:order:" + symbol + ":" + orderId + ":" + action
	cmd := middleware.RedisClient.Exists(key)
	if cmd.Val() == 1 {
		return true
	}
	return false
}

// SavePrice 每个symbol的价格， 如key是matching:price:BTCUSD 则存储对应最新价格
func SavePrice(symbol string, price decimal.Decimal) {
	key := "matching:price:" + symbol
	middleware.RedisClient.Set(key, price.String(), 0)
}

func GetPrice(symbol string) decimal.Decimal {
	key := "matching:price:" + symbol
	priceStr := middleware.RedisClient.Get(key).Val()
	result, err := decimal.NewFromString(priceStr)
	if err != nil {
		result = decimal.Zero
	}
	return result
}

func RemovePrice(symbol string) {
	key := "matching:price:" + symbol
	middleware.RedisClient.Del(key)
}

//SaveOrder
//1.既能缓存下单请求，也能缓存撤单请求；
// 撮合实际在机器本地内存中进行，本机会存到交易账本OrderBook中，服务重启或宕机会丢失OrderBook数据，所以需要缓存到Redis存储
//2.订单要符合定序要求。
// 订单通道里的订单是定序的，交易委托账本里同价格的订单也是按时间排序的，那缓存时如果不定序，程序重启后就难以保证按原有的顺序恢复订单。

// SaveOrder key如何设计
//分两类缓存，第一类保存每个独立的订单请求，包括下单和撤单；第二类分交易标的保存对应 symbol 所有订单请求的订单 ID 和 action。
//第一类， matching:order:{symbol}:{orderId}:{action}，symbol、orderId 和 action 则是对应订单的三个变量值。比如，某订单 symbol = “BTCUSD”，orderId = “12345”，action = “cancel”，那该订单保存到 Redis 的 Key 值就是 matching:order:BTCUSD:12345:cancel。该 Key 对应的 Value 则是保存整个订单对象，可以用 hash 类型存储。
//
//第二类，matching:orderids:{symbol}，Value 保存的是 sorted set 类型的数据，保存对应 symbol 的所有订单请求，每条记录保存的值为 {orderId}:{action}，而 score 值设为对应订单的 {timestamp}。用订单时间作为 score 就可以保证定序了。还记得之前文章我们将订单时间的单位设为 100 纳秒，保证时间戳长度刚好为 16 位吗？这是因为，如果超过 16 位，那 score 将转为科学计数法表示，那将会导致数字失真。
func SaveOrder(order map[string]interface{}) {
	symbol := order["symbol"].(string)
	orderId := order["orderId"].(string)
	timestamp := order["timestamp"].(string)
	action := order["action"].(string)

	key := "matching:order:" + symbol + ":" + orderId + ":" + action
	middleware.RedisClient.HMSet(key, order)

	key = "matching:orderids:" + symbol
	score, _ := strconv.ParseFloat(timestamp, 64)
	z := redis.Z{
		Score:  score,
		Member: orderId + ":" + action,
	}
	middleware.RedisClient.ZAdd(key, z)
}

func HSetOrder(order map[string]interface{}, field string, value interface{}) {
	symbol := order["symbol"].(string)
	orderId := order["orderId"].(string)
	action := order["action"].(string)
	key := "matching:order:" + symbol + ":" + orderId + ":" + action
	middleware.RedisClient.HSet(key, field, value)
}

func RemoveOrder(order map[string]interface{}) {
	symbol := order["symbol"].(string)
	orderId := order["orderId"].(string)
	action := order["action"].(string)
	key := "matching:order:" + symbol + ":" + orderId + ":" + action
	middleware.RedisClient.Del(key)
	key = "matching:orderids:" + symbol
	member := orderId + ":" + action
	middleware.RedisClient.ZRem(key, member)
}

func GetOrderIdsWithAction(symbol string) []string {
	key := "matching:orderids:" + symbol
	return middleware.RedisClient.ZRange(key, 0, -1).Val()
}

func Clear(key string) {

}
