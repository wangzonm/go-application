package process

import (
	"github.com/shopspring/decimal"
	"matching/engine"
	"matching/errcode"
	"matching/middleware/cache"
)

// NewEngine 第一步先判断 ChanMap[symbol] 是否为空，该 ChanMap 就是上文所说的引擎包初始化时用来保存订单通道的 map。如果 ChanMap[symbol] 不为空，说明该 symbol 的撮合引擎已经启动过了，那就返回错误。如果为空，那就初始化这个 symbol 的通道，从代码可知，ChanMap[symbol] 初始化为一个缓冲大小为 100 的订单通道。
//
//接着，就调用 engine.Run() 启动一个 goroutine 了，这行代码即表示用 goroutine 的方式启动指定 symbol 的撮合引擎了。
//
//然后，就将 symbol 和 price 都缓存起来了。
func NewEngine(symbol string, price decimal.Decimal) (book *engine.OrderBookPriority, errCode *errcode.Errcode) {
	if engine.ChanMap[symbol] != nil {
		return book, errcode.EngineExist
	}

	book = &engine.OrderBookPriority{}
	book.Init()
	engine.ChanMap[symbol] = make(chan engine.QueueItem, 10)
	go engine.Run(symbol, price, book)

	cache.SaveSymbol(symbol)
	cache.SavePrice(symbol, price)

	return book, errcode.OK
}
