package types

import (
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/ontio/ontology/common"
)

type TransferBlock struct{
	block.Key
	SingerBlock common.Uint256
	NewOwner common.Uint256
	Coins []PositiveBlock
}

func(transferBlock *TransferBlock) SerializeValue() []byte{
	sink:=common.ZeroCopySink{}
	sink.WriteHash(transferBlock.SingerBlock)
	sink.WriteHash(transferBlock.NewOwner)
	sink.WriteVarUint(uint64(len(transferBlock.Coins)))
	for  _ , positiveBlock := range transferBlock.Coins{
		sink.WriteBytes(positiveBlock.SerializeKey())
		sink.WriteBytes(positiveBlock.SerializeValue())
	}
	return sink.Bytes()
}


