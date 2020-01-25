package types

import (
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/ontio/ontology/common"
)

type PayBlock struct{
	block.Key
	SignerSign common.Uint256
	Coins []PositiveBlock
}

func(payBlock *PayBlock) SerializeValue() [] byte{
	sink:=common.ZeroCopySink{}
	sink.WriteHash(payBlock.SignerSign)
	sink.WriteVarUint(len(payBlock.Coins))
	for _, positiveBlock := range payBlock.Coins{
		sink.WriteBytes(positiveBlock.SerializeKey())
		sink.WriteBytes(positiveBlock.SerializeValue())
	}
	return sink.Bytes()
}
