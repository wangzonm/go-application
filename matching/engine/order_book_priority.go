package engine

import (
	"github.com/shopspring/decimal"
	"sync"
)

type OrderBookPriority struct {
	buyBook       *OrderQueue
	sellBook      *OrderQueue
	priceDigit    int
	quantityDigit int
	w             sync.Mutex
	miniTradeQty  decimal.Decimal
}

func (o *OrderBookPriority) Init() {
	o.priceDigit = 2
	o.quantityDigit = 4
	o.miniTradeQty = decimal.New(1, int32(-o.quantityDigit))
	o.buyBook = NewQueue()
	o.sellBook = NewQueue()
}
