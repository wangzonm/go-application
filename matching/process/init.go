package process

import (
	"matching/engine"
	"matching/enum"
	"matching/middleware/cache"
)

// Init 这一步主要是从缓存加载和恢复各交易标的引擎的启动和所有订单数据。
//1.从缓存读取所有 symbol，即程序重启之前，已经开启了撮合的所有交易标的的 symbol；
//2.从缓存读取每个 symbol 对应的价格，这是程序重启前的最新成交价格；
//3.启动每个 symbol 的撮合引擎；
//4.从缓存读取每个 symbol 的所有订单，这些订单都是按时间顺序排列的；
//5.按顺序将这些订单添加到对应 symbol 的订单通道里去。
func Init() {
	symbols := cache.GetSymbols()
	for _, symbol := range symbols {
		price := cache.GetPrice(symbol)
		NewEngine(symbol, price)
		orderIds := cache.GetOrderIdsWithAction(symbol)
		for _, orderId := range orderIds {
			mapOrder := cache.GetOrder(symbol, orderId)
			order := engine.Order{}
			order.FromMap(mapOrder)
			var item engine.QueueItem
			switch order.Side {
			case enum.SideSell:
				item = engine.NewAskItem(order.PriceType, order.OrderId, order.Price, order.Quantity, order.Amount, order.CreateTime)
			case enum.SideBuy:
				item = engine.NewBidItem(order.PriceType, order.OrderId, order.Price, order.Quantity, order.Amount, order.CreateTime)
			}
			engine.ChanMap[order.Symbol] <- item
		}
	}
}
