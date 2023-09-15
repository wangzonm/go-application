package engine

import (
	"container/list"
	"matching/enum"
)

type IOrderQueueInterface interface {
}

// orderQueue https://cloud.tencent.com/developer/article/1541274
//买单队列和卖单队列可以设计为使用统一的订单队列类型，两者只有价格排序方向不同，那订单队列就可以用一个属性来表示排序方向。队列里的所有订单可以采用二维数组或二维链表来保存，考虑到主要操作是插入和删除，用链表比用数组效率更高。如果想让操作效率更高，那就需要使用更复杂的数据结构了，比如再结合跳表。目前版本为了简单，采用简单的二维链表即可。
//
//使用二维链表的话，那链表中的每个元素保存的就是横向上按时间排序的订单链表，这些订单链表又组成了竖向上按价格排序的链表。
//
//另外，还可以保存一个 Map，将价格作为 Key，将同价格的订单链表作为 Value，这样就能加快同价格订单的查询。

type orderQueue struct {
	sortBy     enum.SortDirection
	parentList *list.List
	elementMap map[string]*list.Element
}

//func (q *orderQueue) init(sortBy enum.SortDirection) {
//	q.sortBy = sortBy
//	q.parentList = list.New()
//	q.elementMap = make(map[string]*list.Element)
//}
//
//func (q *orderQueue) addOrder(order *Order) *list.Element {
//	e := q.parentList.PushBack(order)
//	mapE, ok := q.elementMap[order.Price.String()]
//	if !ok {
//		mapE = list.New().Front()
//		q.elementMap[order.Price.String()] = mapE
//	} else {
//		mapE.Prev()
//	}
//	return e
//}
//
//func (q *orderQueue) getHeadOrder() *list.Element {
//	return q.parentList.Front()
//}
//
//func (q *orderQueue) popHeadOrder() any {
//	e := q.parentList.Front()
//	delete(q.elementMap, e.Value.(*Order).OrderId)
//	return q.parentList.Remove(e)
//}
//
//func (q *orderQueue) removeOrder(e *list.Element) any {
//	delete(q.elementMap, e.Value.(*Order).OrderId)
//	return q.parentList.Remove(e)
//}
//
//func (q *orderQueue) getDepthPrice(depth int) (string, int) {
//	if q.parentList.Len() == 0 {
//		return "", 0
//	}
//	p := q.parentList.Front()
//	i := 1
//	for ; i < depth; i++ {
//		t := p.Next()
//		if t != nil {
//			p = t
//		} else {
//			break
//		}
//	}
//	o := p.Value.(*list.List).Front().Value.(*Order)
//	return o.Price.String(), i
//}

//func dealBuyMarketTop(order *Order, book *orderBook, lastTradePrice *decimal.Decimal, depth int) {
//	priceStr, _ := book.getSellDepthPrice(depth)
//	if priceStr == "" {
//		cancelOrder(order)
//		return
//	}
//	limitPrice, _ := decimal.NewFromString(priceStr)
//LOOP:
//	headOrder := book.getHeadSellOrder()
//	if headOrder != nil && limitPrice.GreaterThanOrEqual(headOrder.Price) {
//		matchTrade(headOrder, order, book, lastTradePrice)
//		if order.Amount.IsPositive() {
//			goto LOOP
//		}
//	} else {
//		cancelOrder(order)
//	}
//}

//func dealBuyMarketOpponent(order *Order, book *orderBook, lastTradePrice *decimal.Decimal) {
//	priceStr, _ := book.getSellDepthPrice(1)
//	if priceStr == "" {
//		cancelOrder(order)
//		return
//	}
//	limitPrice, _ := decimal.NewFromString(priceStr)
//LOOP:
//	headOrder := book.getHeadSellOrder()
//	if headOrder != nil && limitPrice.GreaterThanOrEqual(headOrder.Price) {
//		matchTrade(headOrder, order, book, lastTradePrice)
//		if order.Amount.IsPositive() {
//			goto LOOP
//		}
//	} else {
//		order.Price = limitPrice
//		order.Type = enum.TypeLimit
//		book.addBuyOrder(order)
//		cache.UpdateOrder(order.ToMap())
//		log.Info("engine %s, a order has added to the orderbook: %s", order.Symbol, order.ToJson())
//	}
//}
