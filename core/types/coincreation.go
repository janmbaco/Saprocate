package types

import (
	"github.com/ontio/ontology/common"
)

//Create a coin of from address (in header block) and send to ToAddress
type CoinCreation struct {
	ToAddress common.Address
	Expiration uint64 // in unix time
}

func (this *CoinCreation) Serialization(sink *common.ZeroCopySink) {
	sink.WriteBytes(this.ToAddress[:])
	sink.WriteUint64(this.Expiration)
}

func (this *CoinCreation) Deserialization (source *common.ZeroCopySource)  {
	var buf []byte
	buf, eof := source.NextBytes(common.ADDR_LEN)
	tryEof(eof)
	copy(this.ToAddress[:], buf)
	this.Expiration, eof = source.NextUint64()
	tryEof(eof)
}

func(this *CoinCreation) GetType() BlockType{
	return CoinCreator
}