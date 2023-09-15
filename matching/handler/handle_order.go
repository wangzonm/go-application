package handler

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"io/ioutil"
	"matching/engine"
	"matching/enum"
	"matching/process"
	"net/http"
	"time"
)

type handleOrderParams struct {
	Action    enum.OrderAction `json:"action"`
	Symbol    string           `json:"symbol"`
	OrderId   string           `json:"order_id"`
	Side      enum.OrderSide   `json:"side"`
	Type      enum.OrderType   `json:"type"`
	Amount    decimal.Decimal  `json:"amount"`
	Price     decimal.Decimal  `json:"price"`
	Timestamp int64            `json:"timestamp"`
}

func HandleOrder(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	var params handleOrderParams
	if err := json.Unmarshal(body, &params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	//order := engine.Order{
	//	Action:    params.Action,
	//	Symbol:    params.Symbol,
	//	orderId:   params.OrderId,
	//	Side:      params.Side,
	//	Type:      params.Type,
	//	Amount:    params.Amount,
	//	Price:     params.Price,
	//	Timestamp: params.Timestamp,
	//}
	var order engine.QueueItem
	if params.Side == enum.SideSell {
		quantity := decimal.New(1, 0)
		order = engine.NewAskItem(params.Type, params.OrderId, params.Price, quantity, params.Amount, time.Now().UnixNano()/1e3)
	} else {
		quantity := decimal.New(1, 0)
		order = engine.NewBidItem(params.Type, params.OrderId, params.Price, quantity, params.Amount, time.Now().UnixNano()/1e3)
	}
	order.SetSymbol(params.Symbol)
	errCode := process.Dispatch(order)
	w.Write(errCode.ToJson())
}
