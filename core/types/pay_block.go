package types

import (
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/ontio/ontology/common"
)

type PayBlock struct{
	block.Key
	PrevHash common.Uint256
	From common.Uint256
	Coins []block.Coin
}

func(this *PayBlock) SerializeValue() [] byte{
	sink:= &common.ZeroCopySink{}
	sink.WriteHash(this.From)
	sink.WriteVarUint(uint64(len(this.Coins)))
	for _, coin := range this.Coins{
		coin.Serilize(sink)
	}
	return sink.Bytes()
}
