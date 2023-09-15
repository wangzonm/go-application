package engine

import (
	"fmt"
	"github.com/shopspring/decimal"
	"matching/enum"
	"testing"
	"time"
)

func TestOrderQueue(t *testing.T) {
	q := new(orderQueue)
	q.init(enum.SortDirectionDesc)
	q.addOrder(&Order{
		Action:    enum.ActionCreate,
		Symbol:    "BTCUSDT",
		OrderId:   "BU00001",
		Side:      enum.SideBuy,
		Type:      enum.TypeLimit,
		Amount:    decimal.New(100, 0),
		Price:     decimal.New(100, 0),
		Timestamp: time.Now().Unix(),
	})
	q.addOrder(&Order{
		Action:    enum.ActionCreate,
		Symbol:    "BTCUSDT",
		OrderId:   "BU00002",
		Side:      enum.SideBuy,
		Type:      enum.TypeLimit,
		Amount:    decimal.New(200, 0),
		Price:     decimal.New(200, 0),
		Timestamp: time.Now().Unix(),
	})
	headE := q.popHeadOrder()
	fmt.Printf("headE %#v\n", headE)
	fmt.Printf("list %#v\n", q)
	//deep, num := q.getDepthPrice(2)
	//fmt.Printf("getDepthPrice %#v num = %d\n", deep, num)
}
