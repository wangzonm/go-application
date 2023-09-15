package engine

import (
	"github.com/shopspring/decimal"
	"matching/enum"
)

type Order struct {
	Action     enum.OrderAction `json:"action"`
	Symbol     string           `json:"symbol"`
	Side       enum.OrderSide   `json:"side,string"`
	OrderId    string           `json:"orderId"`
	PriceType  enum.OrderType   `json:"type,string"`
	Amount     decimal.Decimal  `json:"amount"`
	Price      decimal.Decimal  `json:"price"`
	Quantity   decimal.Decimal  `json:"quantity"`
	CreateTime int64            `json:"timestamp,string"`
	index      int
}

func (o *Order) Less(item QueueItem) bool {
	//TODO implement me
	panic("implement me")
}

func (o *Order) GetOrderSide() enum.OrderSide {
	//TODO implement me
	panic("implement me")
}

func (o *Order) GetSymbol() string {
	return o.Symbol
}

func (o *Order) GetAction() enum.OrderAction {
	return o.Action
}

func (o *Order) GetIndex() int {
	return o.index
}

func (o *Order) SetSymbol(symbol string) {
	o.Symbol = symbol
}

func (o *Order) SetAction(actoin enum.OrderAction) {
	o.Action = actoin
}

func (o *Order) SetIndex(index int) {
	o.index = index
}

func (o *Order) SetQuantity(qnt decimal.Decimal) {
	o.Quantity = qnt
}

func (o *Order) SetAmount(amount decimal.Decimal) {
	o.Amount = amount
}

func (o *Order) GetUniqueId() string {
	return o.OrderId
}

func (o *Order) GetPrice() decimal.Decimal {
	return o.Price
}

func (o *Order) GetQuantity() decimal.Decimal {
	return o.Quantity
}

func (o *Order) GetCreateTime() int64 {
	return o.CreateTime
}

func (o *Order) GetPriceType() enum.OrderType {
	return o.PriceType
}
func (o *Order) GetAmount() decimal.Decimal {
	return o.Amount
}

// 这个方法留在具体的 ask/bid 队列中实现
// func (o *Order) Less() {}

type AskItem struct {
	Order
}

func (a *AskItem) Less(other QueueItem) bool {
	//价格优先，时间优先原则
	//价格低的在最上面
	return (a.Price.Cmp(other.(*AskItem).Price) == -1) || (a.Price.Cmp(other.(*AskItem).Price) == 0 && a.CreateTime < other.(*AskItem).CreateTime)
}

func (a *AskItem) GetOrderSide() enum.OrderSide {
	return enum.SideSell
}

type BidItem struct {
	Order
}

func (a *BidItem) Less(other QueueItem) bool {
	//价格优先，时间优先原则
	//价格高的在最上面
	return (a.Price.Cmp(other.(*BidItem).Price) == 1) || (a.Price.Cmp(other.(*BidItem).Price) == 0 && a.CreateTime < other.(*BidItem).CreateTime)
}

func (a *BidItem) GetOrderSide() enum.OrderSide {
	return enum.SideBuy
}

func NewAskItem(pt enum.OrderType, uniqId string, price, quantity, amount decimal.Decimal, createTime int64) *AskItem {
	return &AskItem{
		Order: Order{
			Symbol:     "BTC-USDT",
			Action:     enum.ActionCreate,
			Side:       enum.SideSell,
			OrderId:    uniqId,
			Price:      price,
			Quantity:   quantity,
			CreateTime: createTime,
			PriceType:  pt,
			Amount:     amount,
		},
	}
}

func NewAskLimitItem(uniq string, price, quantity decimal.Decimal, createTime int64) *AskItem {
	return NewAskItem(enum.TypeLimit, uniq, price, quantity, decimal.Zero, createTime)
}

func NewAskMarketQtyItem(uniq string, quantity decimal.Decimal, createTime int64) *AskItem {
	return NewAskItem(enum.PriceTypeMarketQuantity, uniq, decimal.Zero, quantity, decimal.Zero, createTime)
}

//市价 按金额卖出订单时，需要用户持有交易物的数量，在撮合时候防止超卖
func NewAskMarketAmountItem(uniq string, amount, maxHoldQty decimal.Decimal, createTime int64) *AskItem {
	return NewAskItem(enum.PriceTypeMarketAmount, uniq, decimal.Zero, maxHoldQty, amount, createTime)
}

func NewBidItem(pt enum.OrderType, uniqId string, price, quantity, amount decimal.Decimal, createTime int64) *BidItem {
	return &BidItem{
		Order: Order{
			Symbol:     "BTC-USDT",
			Action:     enum.ActionCreate,
			Side:       enum.SideBuy,
			OrderId:    uniqId,
			Price:      price,
			Quantity:   quantity,
			CreateTime: createTime,
			PriceType:  pt,
			Amount:     amount,
		}}
}

func NewBidLimitItem(uniq string, price, quantity decimal.Decimal, createTime int64) *BidItem {
	return NewBidItem(enum.TypeLimit, uniq, price, quantity, decimal.Zero, createTime)
}

func NewBidMarketQtyItem(uniq string, quantity, maxAmount decimal.Decimal, createTime int64) *BidItem {
	return NewBidItem(enum.PriceTypeMarketQuantity, uniq, decimal.Zero, quantity, maxAmount, createTime)
}

func NewBidMarketAmountItem(uniq string, amount decimal.Decimal, createTime int64) *BidItem {
	return NewBidItem(enum.PriceTypeMarketAmount, uniq, decimal.Zero, decimal.Zero, amount, createTime)
}
