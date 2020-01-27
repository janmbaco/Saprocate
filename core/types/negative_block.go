package types

import (
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/ontio/ontology/common"
)

type NegativeBlock struct{
	block.Key
	Previous *block.Key
	PositiveBlock *block.Key
}

func(this *NegativeBlock) SerializeValue() []byte{
	sink := &common.ZeroCopySink{}
	this.Previous.Serialize(sink)
	this.PositiveBlock.Serialize(sink)
	return sink.Bytes()
}

