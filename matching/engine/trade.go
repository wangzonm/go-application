package engine

type Trade struct {
	MakerId   int64
	TakerId   int64
	TakerSide int64
	Amount    int64
	Price     int64
	Timestamp int64
}

//func matchBuyTrade(sellOrder, buyOrder *Order, book *orderBook, lastTradePrice *decimal.Decimal) {
//	log.Info("matchTrade_sell_before buy_order:%+v buy_order:%+v", buyOrder, sellOrder)
//	buyOrder.Amount = buyOrder.Amount.Sub(sellOrder.Amount)
//	buyOrder.Price = buyOrder.Price.Sub(sellOrder.Price)
//	book.sellBook.removeOrder(book.sellBook.elementMap[sellOrder.OrderId])
//	if buyOrder.Amount.IsZero() {
//		book.buyBook.removeOrder(book.buyBook.elementMap[buyOrder.OrderId])
//	}
//	log.Info("matchTrade_sell_after sell_order:%+v buy_order:%+v", sellOrder, buyOrder)
//}
//
//func matchSellTrade(buyOrder, sellOrder *Order, book *orderBook, lastTradePrice *decimal.Decimal) {
//	log.Info("matchTrade_sell_before sell_order:%+v buy_order:%+v", sellOrder, buyOrder)
//	sellOrder.Amount = sellOrder.Amount.Sub(buyOrder.Amount)
//	sellOrder.Price = sellOrder.Price.Sub(buyOrder.Price)
//	//TODO is pop safe ?
//	book.buyBook.popHeadOrder()
//	if sellOrder.Amount.IsZero() {
//		book.sellBook.removeOrder(book.buyBook.elementMap[buyOrder.OrderId])
//	}
//	log.Info("matchTrade_sell_after sell_order:%+v buy_order:%+v", sellOrder, buyOrder)
//}
