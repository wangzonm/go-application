package engine

import (
	"github.com/shopspring/decimal"
	"matching/enum"
	"matching/log"
	"matching/middleware/cache"
	"time"
)

var Debug = false
var Count = 0

//Run 第一步，先定义和初始化了一个 book 变量，该变量就是用来保存整个交易委托账本。
//
//接着，就是一个 for 循环了，for 循环里的第一行就是从对应 symbol 的订单通道里读取出一个订单，读取到订单时，order 变量就会有值，且 ok 变量为 true。如果通道里暂时没有订单，那就会阻塞在这行代码，直到从通道中获取到订单或通道已关闭的消息。
//
//当通道被关闭之后，最后，从通道中读取到的 ok 变量则为 false，当然，在这之前，会先依序读取完通道里剩下的订单。当 ok 为 false 时，引擎里会执行两步操作：一是从 ChanMap 中删除该 symbol 对应的记录，二是清空该 symbol 对应的缓存数据。最后用 return 来退出 for 循环，这样，整个 Run() 函数就结束退出了，意味着该引擎也真正关闭了。
//
//当每读取到一个订单，就会判断是下单还是撤单，然后进行相应的逻辑处理了。
func Run(symbol string, price decimal.Decimal, book *OrderBookPriority) {
	lastTradePrice := price
	log.Info("engine %s is running", symbol)
	for {
		select {
		case order := <-ChanMap[symbol]:
			log.Info("engine %s receive an order: %s", symbol, order.ToJson())
			switch order.GetAction() {
			case enum.ActionCreate:
				dealCreate(order, book, &lastTradePrice)
			case enum.ActionCancel:
				//dealCancel(&order, book)
			}
		default:
			book.handlerLimitOrder()
		}
		//order, ok := <-ChanMap[symbol]
		//if !ok {
		//	log.Info("engine %s is closed", symbol)
		//	delete(ChanMap, symbol)
		//	cache.RemoveSymbol(symbol)
		//	return
		//}
		//log.Info("engine %s receive an order: %s", symbol, order.ToJson())
		//switch order.GetAction() {
		//case enum.ActionCreate:
		//	dealCreate(order, book, &lastTradePrice)
		//case enum.ActionCancel:
		//	//dealCancel(&order, book)
		//}
	}
}

func dealCreate(order QueueItem, book *OrderBookPriority, lastTradePrice *decimal.Decimal) {
	switch order.GetOrderSide() {
	case enum.SideBuy:
		book.buyBook.Push(order)
	case enum.SideSell:
		book.sellBook.Push(order)
	}
	switch order.GetPriceType() {
	case enum.TypeLimit:
		book.handlerLimitOrder()
		//dealLimit(order, book, lastTradePrice)
	case enum.TypeLimitIoc:
		//dealLimitIoc(order, book, lastTradePrice)
	default:
		dealMarket(order, book, lastTradePrice)
		//case enum.TypeMarketTop5:
		//	dealMarketTop5(order, book, lastTradePrice)
		//case enum.TypeMarketTop10:
		//	dealMarketTop10(order, book, lastTradePrice)
		//case enum.TypeMarketOpponent:
		//	dealMarketOpponent(order, book, lastTradePrice)
	}
}

// dealCancel 1.从委托账本中移除该订单；2.从缓存中移除该订单；3.发送撤单结果到 MQ。
//func dealCancel(order *Order, book *orderBook) {
//	var ok bool
//	switch order.Side {
//	case enum.SideBuy:
//		ok = book.removeBuyOrder(order)
//	case enum.SideSell:
//		ok = book.removeSellOrder(order)
//	}
//
//	cache.RemoveOrder(order.ToMap())
//	//mq.SendCancelResult(order.Symbol, order.OrderId, ok)
//	log.Info("engine %s, order %s cancel result is %s", order.Symbol, order.GetUniqueId(), ok)
//}

//func dealLimit(order QueueItem, book *OrderBookPriority, lastTradePrice *decimal.Decimal) {
//	switch order.GetOrderSide() {
//	case enum.SideBuy:
//		book.buyBook.Push(order)
//	case enum.SideSell:
//		book.sellBook.Push(order)
//	}
//}

func dealLimitIoc(order *Order, book *orderBook, lastTradePrice *decimal.Decimal) {

}

func (t *OrderBookPriority) handlerLimitOrder() {
	ok := func() bool {
		t.w.Lock()
		defer t.w.Unlock()

		if t.sellBook == nil || t.buyBook == nil {
			return false
		}

		if t.sellBook.Len() == 0 || t.buyBook.Len() == 0 {
			return false
		}

		askTop := t.sellBook.Top()
		bidTop := t.buyBook.Top()

		defer func() {
			if askTop.GetQuantity().Equal(decimal.Zero) {
				t.sellBook.Remove(askTop.GetUniqueId())
			}
			if bidTop.GetQuantity().Equal(decimal.Zero) {
				t.buyBook.Remove(bidTop.GetUniqueId())
			}
		}()

		if bidTop.GetPrice().Cmp(askTop.GetPrice()) >= 0 {
			curTradeQty := decimal.Zero
			curTradePrice := decimal.Zero
			if bidTop.GetQuantity().Cmp(askTop.GetQuantity()) >= 0 {
				curTradeQty = askTop.GetQuantity()
			} else if bidTop.GetQuantity().Cmp(askTop.GetQuantity()) == -1 {
				curTradeQty = bidTop.GetQuantity()
			}
			askTop.SetQuantity(askTop.GetQuantity().Sub(curTradeQty))
			bidTop.SetQuantity(bidTop.GetQuantity().Sub(curTradeQty))

			if askTop.GetCreateTime() >= bidTop.GetCreateTime() {
				curTradePrice = bidTop.GetPrice()
			} else {
				curTradePrice = askTop.GetPrice()
			}

			t.sendTradeResultNotify(askTop, bidTop, curTradePrice, curTradeQty, "")
			return true
		} else {
			return false
		}

	}()

	if !ok {
		time.Sleep(time.Duration(60) * time.Millisecond)
	} else {
		if Debug {
			time.Sleep(time.Second * time.Duration(1))
		}
	}
}

func dealMarket(order QueueItem, t *OrderBookPriority, lastTradePrice *decimal.Decimal) {
	t.w.Lock()
	defer t.w.Unlock()
	if order.GetOrderSide() == enum.SideSell {
		//t.doMarketSell(order)
	} else {
		t.doBorrow(order)
	}
}

func (t *OrderBookPriority) doBorrow(item QueueItem) {
	for {
		ok := func() bool {
			if t.sellBook.Len() == 0 {
				return false
			}

			ask := t.sellBook.Top()
			if item.GetAmount().Cmp(ask.GetAmount()) > 0 {
				return false
			}

			ask.SetAmount(ask.GetAmount().Sub(item.GetAmount()))
			t.sendMatchResult(ask, item, item.GetAmount().String(), item.GetUniqueId())
			return false
		}()
		if !ok {
			break
		}
	}
}

func (t *OrderBookPriority) sendMatchResult(ask, bid QueueItem, price, market_done string) {
	Count++
	log.Info("ask=%+v bid=%+v Price=%v market_done=%v count:%v", ask, bid, price, market_done, Count)
	cache.RemoveOrder(map[string]interface{}{
		"symbol":  bid.GetSymbol(),
		"orderId": bid.GetUniqueId(),
		"action":  bid.GetAction().String(),
	})
	order := ask.ToMap()
	cache.HSetOrder(order, "amount", order["amount"])
}

func (t *OrderBookPriority) doMarketSell(item QueueItem) {

	for {
		ok := func() bool {

			if t.buyBook.Len() == 0 {
				return false
			}

			bid := t.buyBook.Top()
			if item.GetPriceType() == enum.PriceTypeMarketQuantity {

				curTradeQuantity := decimal.Zero
				//市价按买入数量
				if item.GetQuantity().Equal(decimal.Zero) {
					return false
				}

				if bid.GetQuantity().Cmp(item.GetQuantity()) <= 0 {
					curTradeQuantity = bid.GetQuantity()
					t.buyBook.Remove(bid.GetUniqueId())
				} else {
					curTradeQuantity = item.GetQuantity()
					bid.SetQuantity(bid.GetQuantity().Sub(curTradeQuantity))
				}

				item.SetQuantity(item.GetQuantity().Sub(curTradeQuantity))

				//退出条件
				// a.对面订单空了
				// b.市价订单完全成交了
				if t.buyBook.Len() == 0 || item.GetQuantity().Equal(decimal.Zero) {
					t.sendTradeResultNotify(item, bid, bid.GetPrice(), curTradeQuantity, item.GetUniqueId())
				} else {
					t.sendTradeResultNotify(item, bid, bid.GetPrice(), curTradeQuantity, "")
				}

				return true
			} else if item.GetPriceType() == enum.PriceTypeMarketAmount {
				//市价-按成交金额成交
				if bid.GetPrice().Cmp(decimal.Zero) <= 0 {
					return false
				}

				maxQty := func(amount, price, needQty decimal.Decimal) decimal.Decimal {
					a := amount.Div(price)
					a = a.Truncate(int32(t.quantityDigit))
					return decimal.Min(a, needQty)
				}
				maxTradeQty := maxQty(item.GetAmount(), bid.GetPrice(), item.GetQuantity())

				curTradeQty := decimal.Zero
				if maxTradeQty.Cmp(t.miniTradeQty) < 0 {
					return false
				}

				if bid.GetQuantity().Cmp(maxTradeQty) <= 0 {
					curTradeQty = bid.GetQuantity()
					t.buyBook.Remove(bid.GetUniqueId())
				} else {
					curTradeQty = maxTradeQty
					bid.SetQuantity(bid.GetQuantity().Sub(curTradeQty))
				}

				item.SetAmount(item.GetAmount().Sub(curTradeQty.Mul(bid.GetPrice())))
				//市价 按成交额卖出时，需要用户持有的资产数量来进行限制
				item.SetQuantity(item.GetQuantity().Sub(curTradeQty))

				//退出条件
				// a.对面订单空了
				// b.金额完全成交
				// c.剩余资金不满足最小成交量
				if t.buyBook.Len() == 0 || maxQty(item.GetAmount(), t.buyBook.Top().GetPrice(), item.GetQuantity()).Cmp(t.miniTradeQty) < 0 {
					t.sendTradeResultNotify(item, bid, bid.GetPrice(), curTradeQty, item.GetUniqueId())
				} else {
					t.sendTradeResultNotify(item, bid, bid.GetPrice(), curTradeQty, "")
				}

				return true
			}

			return false
		}()

		if !ok {
			break
		}

	}
}

func (t *OrderBookPriority) doMarketBuy(item QueueItem) {

	for {
		ok := func() bool {

			if t.sellBook.Len() == 0 {
				return false
			}

			ask := t.sellBook.Top()
			if item.GetPriceType() == enum.PriceTypeMarketQuantity {
				maxQty := func(remainAmount, marketPrice, needQty decimal.Decimal) decimal.Decimal {
					qty := remainAmount.Div(marketPrice)
					return decimal.Min(qty, needQty)
				}
				maxTradeQty := maxQty(item.GetAmount(), ask.GetPrice(), item.GetQuantity())
				curTradeQty := decimal.Zero

				//市价按买入数量
				if maxTradeQty.Cmp(t.miniTradeQty) < 0 {
					return false
				}

				if ask.GetQuantity().Cmp(maxTradeQty) <= 0 {
					curTradeQty = ask.GetQuantity()
					t.sellBook.Remove(ask.GetUniqueId())
				} else {
					curTradeQty = maxTradeQty
					ask.SetQuantity(ask.GetQuantity().Sub(curTradeQty))
				}

				item.SetQuantity(item.GetQuantity().Sub(curTradeQty))
				item.SetAmount(item.GetAmount().Sub(curTradeQty.Mul(ask.GetPrice())))

				//检查本次循环撮合是否是该订单最后一次撮合
				//如果是则标记该市价订单已经完成了
				//结束的条件：
				// a.对面订单列表空了
				// b.已经达到了用户需要的数量
				// c.剩余资金已经不能达到最小成交需求
				if t.sellBook.Len() == 0 || item.GetQuantity().Equal(decimal.Zero) ||
					maxQty(item.GetAmount(), t.sellBook.Top().GetPrice(), item.GetQuantity()).Cmp(t.miniTradeQty) < 0 {
					t.sendTradeResultNotify(ask, item, ask.GetPrice(), curTradeQty, item.GetUniqueId())
				} else {
					t.sendTradeResultNotify(ask, item, ask.GetPrice(), curTradeQty, "")
				}

				return true
			} else if item.GetPriceType() == enum.PriceTypeMarketAmount {
				//市价-按成交金额
				//成交金额不包含手续费，手续费应该由上层系统计算提前预留
				//撮合会针对这个金额最大限度的买入
				if ask.GetPrice().Cmp(decimal.Zero) <= 0 {
					return false
				}

				maxQty := func(amount, price decimal.Decimal) decimal.Decimal {
					return amount.Div(price)
				}

				maxTradeQty := maxQty(item.GetAmount(), ask.GetPrice()) //item.GetAmount().Div(ask.GetPrice())
				curTradeQty := decimal.Zero

				if maxTradeQty.Cmp(t.miniTradeQty) < 0 {
					return false
				}
				if ask.GetQuantity().Cmp(maxTradeQty) <= 0 {
					curTradeQty = ask.GetQuantity()
					t.sellBook.Remove(ask.GetUniqueId())
				} else {
					curTradeQty = maxTradeQty
					ask.SetQuantity(ask.GetQuantity().Sub(curTradeQty))
				}

				//部分成交了，需要更新这个单的剩余可成交金额，用于下一轮重新计算最大成交量
				item.SetAmount(item.GetAmount().Sub(curTradeQty.Mul(ask.GetPrice())))
				item.SetQuantity(item.GetQuantity().Add(curTradeQty))

				//检查本次循环撮合是否是该订单最后一次撮合
				//结束的条件：
				// a.对面订单列表空了
				// b.已经达到了用户需要的数量
				// c.剩余资金已经不能达到最小成交需求
				if t.sellBook.Len() == 0 || item.GetQuantity().Equal(decimal.Zero) ||
					maxQty(item.GetAmount(), t.sellBook.Top().GetPrice()).Cmp(t.miniTradeQty) < 0 {
					t.sendTradeResultNotify(ask, item, ask.GetPrice(), curTradeQty, item.GetUniqueId())
				} else {
					t.sendTradeResultNotify(ask, item, ask.GetPrice(), curTradeQty, "")
				}
				return true
			}

			return false
		}()

		if !ok {
			break
		}

	}
}

func (t *OrderBookPriority) sendTradeResultNotify(ask, bid QueueItem, price, tradeQty decimal.Decimal, market_done string) {
	Count++
	log.Info("ask=%+v bid=%+v Price=%v tradeQty=%+v market_done=%v count:%v", ask, bid, price, tradeQty, market_done, Count)

	if market_done == "" {
		cache.RemoveOrder(map[string]interface{}{
			"symbol":  ask.GetSymbol(),
			"orderId": ask.GetUniqueId(),
			"action":  ask.GetAction().String(),
		})
		cache.RemoveOrder(map[string]interface{}{
			"symbol":  bid.GetSymbol(),
			"orderId": bid.GetUniqueId(),
			"action":  bid.GetAction().String(),
		})
	}
}

//dealBuyLimit
//1.从委托账本中读取出卖单队列的头部订单；
//2.如果头部订单为空，或新订单(买单)价格小于头部订单(卖单)，则无法匹配成交，那就将新订单添加到委托账本的买单队列中去；
//3.如果头部订单不为空，且新订单(买单)价格大于等于头部订单(卖单)，则两个订单可以匹配成交，那就对这两个订单进行成交处理；
//4.如果上一步的成交处理完之后，新订单的剩余数量还不为零，那就继续重复第一步。
//
//其中，匹配成交的记录会作为一条输出记录发送到 MQ。
//func dealBuyLimit(order *Order, book *orderBook, lastTradePrice *decimal.Decimal) {
//LOOP:
//	headOrder := book.getHeadSellOrder()
//	if headOrder == nil || order.Price.LessThan(headOrder.Price) {
//		book.addBuyOrder(order)
//		log.Info("engine %s, a order has added to the buy_orderbook: %s", order.Symbol, order.ToJson())
//	} else {
//		matchBuyTrade(headOrder, order, book, lastTradePrice)
//		if order.Amount.IsPositive() {
//			goto LOOP
//		}
//	}
//}

//func dealSellLimit(order *Order, book *orderBook, lastTradePrice *decimal.Decimal) {
//LOOP:
//	headOrder := book.getHeadBuyOrder()
//	if headOrder == nil || order.Price.LessThan(headOrder.Price) {
//		book.addSellOrder(order)
//		log.Info("engine %s, a order has added to the sell_orderbook: %s", order.Symbol, order.ToJson())
//	} else {
//		matchSellTrade(headOrder, order, book, lastTradePrice)
//		if order.Amount.IsPositive() {
//			goto LOOP
//		}
//	}
//}
