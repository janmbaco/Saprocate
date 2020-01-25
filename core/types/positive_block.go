package types

import (
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/ontio/ontology/common"
)

type PositiveBlock struct{
	block.Key
	Numus common.Uint256
	Timestamp uint64
}

func(positiveBlock *PositiveBlock) SerializeValue() []byte{
	sink := common.ZeroCopySink{}
	sink.WriteHash(positiveBlock.Numus)
	sink.WriteUint64(positiveBlock.Timestamp)
	return sink.Bytes()
}
