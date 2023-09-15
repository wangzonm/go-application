package process

import (
	"matching/engine"
	"matching/enum"
	"matching/errcode"
	"matching/middleware/cache"
)

// Dispatch 第一步，判断 ChanMap[order.Symbol] 是否为空，如果为空，表示引擎没开启，那就无法处理订单。
//
//第二步，判断订单是否存在。如果是 create 订单，那缓存中就不应该查到订单，否则说明是重复请求。如果是 cancel 订单，那缓存中如果也查不到订单，那说明该订单已经全部成交或已经成功撤单过了。
//
//第三步，将订单时间设为当前时间，时间单位是 100 纳秒，这可以保证时间戳长度刚好为 16 位，保存到 Redis 里就不会有精度失真的问题。这点后续文章讲到 Redis 详细设计时再说。
//
//第四步，将订单缓存。
//
//第五步，将订单传入对应的订单通道，对应引擎会从该通道中获取该订单进行处理。这一步就实现了订单的分发。
//
//第六步，返回 OK。
func Dispatch(order engine.QueueItem) *errcode.Errcode {
	if engine.ChanMap[order.GetSymbol()] == nil {
		return errcode.EngineNotFound
	}

	if order.GetAction() == enum.ActionCreate {
		if cache.OrderExist(order.GetSymbol(), order.GetUniqueId(), order.GetAction().String()) {
			return errcode.OrderExist
		}
	} else {
		if !cache.OrderExist(order.GetSymbol(), order.GetUniqueId(), enum.ActionCreate.String()) {
			return errcode.OrderNotFound
		}
	}

	cache.SaveOrder(order.ToMap())
	engine.ChanMap[order.GetSymbol()] <- order

	return errcode.OK
}
