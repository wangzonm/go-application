package engine

type IOrderBook interface {
	init()
	addBuyOrder(order *Order)
	addSellOrder(order *Order)
	getHeadBuyOrder() *Order
	getHeadSellOrder() *Order
	popHeadBuyOrder()
	popHeadSellOrder()
	removeHeadBuyOrder(order *Order)
	removeHeadSellOrder(order *Order)
}

//
//  orderBook
//  @Description: https://cloud.tencent.com/developer/article/1550780
//交易委托账本其实就是由两个订单队列组成的，一个买单队列，一个卖单队列。
//
//sortBy 指定价格排序的方向，买单队列是降序的，而卖单队列则是升序的。
//parentList 保存整个二维链表的所有订单，第一维以价格排序，第二维以时间排序。
//elementMap 则是 Key 为价格、Value 为第二维订单链表的键值对。
//
type orderBook struct {
	buyBook  *orderQueue
	sellBook *orderQueue
}

//func (o *orderBook) init() {
//	o.buyBook = &orderQueue{}
//	o.buyBook.init(enum.SortDirectionDesc)
//	o.sellBook = &orderQueue{}
//	o.sellBook.init(enum.SortDirectionAsc)
//}
//
//func (o *orderBook) sortBook() {
//
//}
//
//func (o *orderBook) addBuyOrder(order *Order) {
//	o.buyBook.addOrder(order)
//}

//func (o *orderBook) addSellOrder(order *Order) {
//	order.Timestamp = time.Now().UnixNano() / 1e3
//	o.sellBook.addOrder(order)
//}
//
//func (o *orderBook) removeBuyOrder(order *Order) bool {
//	e := o.buyBook.elementMap[order.OrderId]
//	res := o.buyBook.removeOrder(e)
//	return res.(bool)
//}
//
//func (o *orderBook) removeSellOrder(order *Order) bool {
//	e := o.sellBook.elementMap[order.OrderId]
//	res := o.sellBook.removeOrder(e)
//	return res.(bool)
//}

//func (o *orderBook) getHeadBuyOrder() *Order {
//	if o.buyBook == nil {
//		return nil
//	}
//
//	e := o.buyBook.getHeadOrder()
//	return e.Value.(*Order)
//}
//
//func (o *orderBook) getHeadSellOrder() *Order {
//	if o.sellBook == nil {
//		return nil
//	}
//
//	e := o.sellBook.getHeadOrder()
//	if e == nil {
//		return nil
//	}
//
//	return e.Value.(*Order)
//}
//
//func (o *orderBook) popHeadBuyOrder() {
//	o.buyBook.popHeadOrder()
//}
//
//func (o *orderBook) popHeadSellOrder() {
//	o.sellBook.popHeadOrder()
//}
//
//func (o *orderBook) removeHeadBuyOrder(order *Order) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (o *orderBook) removeHeadSellOrder(order *Order) {
//	//TODO implement me
//	panic("implement me")
//}
