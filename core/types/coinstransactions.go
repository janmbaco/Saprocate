package types

import (
	"github.com/ontio/ontology/common"
)

//Send a transaction coin from address (in header) to address
type CoinsTransaction struct {
	ToAddress common.Address //address where the transaction coins go
	CoinsMove []common.Uint256 //before transaction block hash, only accept type 1 y 2
}

func (this *CoinsTransaction) Serialization(sink *common.ZeroCopySink) {
	sink.WriteBytes(this.ToAddress[:])
	sink.WriteVarUint(uint64(len(this.CoinsMove)))

	for _, txhash := range this.CoinsMove{
		sink.WriteHash(txhash)
	}
}

func (this *CoinsTransaction) Deserialization (source *common.ZeroCopySource)  {
	var buf []byte
	buf, eof := source.NextBytes(common.ADDR_LEN)
	tryEof(eof)
	copy(this.ToAddress[:], buf)
	txlen, eof := source.NextUint64()
	tryEof(eof)
	for i := uint64(0); i < txlen; i++{
		move, eof := source.NextHash()
		tryEof(eof)
		this.CoinsMove = append(this.CoinsMove, move)
	}
}

func(this *CoinsTransaction) GetType() BlockType{
	return CoinsTransactor
}