package types

import (
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/ontio/ontology/common"
)

type PositiveBlock struct{
	block.Key
	Previous *block.Key
	Coin *block.Coin
}

func(this *PositiveBlock) SerializeValue() []byte{
	sink := &common.ZeroCopySink{}
	this.Previous.Serialize(sink)
	this.Coin.Serilize(sink)
	return sink.Bytes()
}
