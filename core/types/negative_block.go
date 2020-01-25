package types

import (
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/ontio/ontology/common"
)

type NegativeBlock struct{
	block.Key
	PositiveBlockSing common.Uint256
	SignerBlock common.Uint256
}

func(negativeBlock *NegativeBlock) SerializeValue() []byte{
	sink:=common.ZeroCopySink{}
	sink.WriteHash(negativeBlock.PositiveBlockSing)
	sink.WriteHash(negativeBlock.SignerBlock)
	return sink.Bytes()
}

