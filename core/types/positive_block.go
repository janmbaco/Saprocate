package types

import (
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/ontio/ontology/common"
)

type PositiveBlock struct{
	block.Key
	PrevHash common.Uint256
	Coin block.Coin
}

func(this *PositiveBlock) SerializeValue() []byte{
	sink := &common.ZeroCopySink{}
	sink.WriteHash(this.PrevHash)
	this.Coin.Serilize(sink)
	return sink.Bytes()
}
