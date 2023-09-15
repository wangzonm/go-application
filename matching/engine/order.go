package engine

import (
	"encoding/json"
)

//
//  Order
//  @Description: 委托单
//
//type Order struct {
//	Action    enum.OrderAction `json:"action"`
//	Symbol    string           `json:"symbol"`
//	OrderId   string           `json:"OrderId"`
//	Side      enum.OrderSide   `json:"side,string"`
//	Type      enum.OrderType   `json:"type,string"`
//	Amount    decimal.Decimal  `json:"Amount"`
//	Price     decimal.Decimal  `json:"Price"`
//	Timestamp int64            `json:"timestamp,string"`
//}

func (o *Order) FromMap(mapOrder map[string]string) {
	b, _ := json.Marshal(&mapOrder)
	_ = json.Unmarshal(b, o)
}

func (o *Order) ToMap() map[string]interface{} {
	m := make(map[string]interface{})
	b, _ := json.Marshal(o)
	_ = json.Unmarshal(b, &m)
	return m
}

func (o *Order) ToJson() string {
	b, _ := json.Marshal(o)
	return string(b)
}
