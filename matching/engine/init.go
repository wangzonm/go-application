package engine

var ChanMap map[string]chan QueueItem

func Init() {
	// 初始化了一个 map，用来保存不同交易标的的订单 channel，作为各交易标的的定序队列来用
	ChanMap = make(map[string]chan QueueItem)
}
