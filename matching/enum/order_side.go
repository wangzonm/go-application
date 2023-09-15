package enum

type OrderSide int

const (
	SideBuy OrderSide = iota + 1
	SideSell
)
