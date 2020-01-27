package types

import (
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/ontio/ontology/common"
)

type TransferBlock struct{
	block.Key
	Previous *block.Key
	From *block.Key
	To *block.Key
	Coins []*block.Coin
}

func(this *TransferBlock) SerializeValue() []byte{
	sink:=&common.ZeroCopySink{}
	this.Previous.Serialize(sink)
	this.From.Serialize(sink)
	this.To.Serialize(sink)
	sink.WriteVarUint(uint64(len(this.Coins)))
	for  _ , coin := range this.Coins{
		coin.Serilize(sink)
	}
	return sink.Bytes()
}


