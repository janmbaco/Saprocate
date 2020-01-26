package types

import (
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/ontio/ontology/common"
)

type NegativeBlock struct{
	block.Key
	PrevHash common.Uint256
	PositiveBlockSing common.Uint256
	SignerBlock common.Uint256
}

func(this *NegativeBlock) SerializeValue() []byte{
	sink:=common.ZeroCopySink{}
	sink.WriteHash(this.PositiveBlockSing)
	sink.WriteHash(this.SignerBlock)
	return sink.Bytes()
}

