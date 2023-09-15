package process

import (
	"matching/engine"
	"matching/errcode"
)

func CloseEngine(symbol string) *errcode.Errcode {
	if engine.ChanMap[symbol] == nil {
		return errcode.EngineNotFound
	}

	close(engine.ChanMap[symbol])

	return errcode.OK
}
