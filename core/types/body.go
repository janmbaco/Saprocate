package types

import (
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/ontio/ontology/common"
	"io"
)

type BlockType uint8

const(
	CurrencyCreator BlockType = iota
	CoinCreator
	CoinsTransactor
)

type Body interface {
	Serialization(sink *common.ZeroCopySink)
	Deserialization(source *common.ZeroCopySource)
	GetType() BlockType
}


func tryEof(eof bool) {
	if eof {
		cross.TryPanic(io.ErrUnexpectedEOF)
	}
}






