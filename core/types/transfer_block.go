package types

import (
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/ontio/ontology/common"
)

type TransferBlock struct{
	block.Key
	PrevHash common.Uint256
	From common.Uint256
	To common.Uint256
	Coins []block.Coin
}

func(this *TransferBlock) SerializeValue() []byte{
	sink:=&common.ZeroCopySink{}
	sink.WriteHash(this.PrevHash)
	sink.WriteHash(this.From)
	sink.WriteHash(this.To)
	sink.WriteVarUint(uint64(len(this.Coins)))
	for  _ , coin := range this.Coins{
		coin.Serilize(sink)
	}
	return sink.Bytes()
}


