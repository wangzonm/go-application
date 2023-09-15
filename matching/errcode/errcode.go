package errcode

import "encoding/json"

type Errcode struct {
	Code int
	Msg  string
}

var (
	OK             = New(200, "ok")
	EngineExist    = New(401, "EngineExist")
	EngineNotFound = New(404, "EngineNotFound")
	BlankSymbol    = New(401, "BlankSymbol")
	InvalidPrice   = New(401, "InvalidPrice")
	OrderExist     = New(401, "OrderExist")
	OrderNotFound  = New(404, "OrderNotFound")
)

func New(code int, msg string) *Errcode {
	return &Errcode{
		Code: code,
		Msg:  msg,
	}
}

func (e *Errcode) ToJson() []byte {
	b, _ := json.Marshal(e)
	return b
}

func (e *Errcode) IsOK() bool {
	if e.Code != 200 {
		return false
	}
	return true
}
