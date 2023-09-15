package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"

	"matching/errcode"
	"matching/process"

	"github.com/shopspring/decimal"
)

type openMatchingParams struct {
	Symbol string          `json:"symbol"`
	Price  decimal.Decimal `json:"price"`
}

func OpenMatching(w http.ResponseWriter, r *http.Request) {
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

	var params openMatchingParams
	if err := json.Unmarshal(body, &params); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if strings.TrimSpace(params.Symbol) == "" {
		w.Write(errcode.BlankSymbol.ToJson())
		return
	}

	if params.Price.IsNegative() {
		w.Write(errcode.InvalidPrice.ToJson())
		return
	}

	if _, e := process.NewEngine(params.Symbol, params.Price); !e.IsOK() {
		w.Write(e.ToJson())
		return
	}

	w.Write(errcode.OK.ToJson())
}

func CloseMatching(w http.ResponseWriter, r *http.Request) {

}
