package enum

type OrderType int

const (
	TypeLimit OrderType = iota + 1
	TypeLimitIoc
	TypeMarket
	TypeMarketTop5
	TypeMarketTop10
	TypeMarketOpponent
	PriceTypeMarketQuantity
	PriceTypeMarketAmount
)
