package header

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/ontio/ontology/common"
)

type Key struct{
	Type blockpkg.Type
	Hash common.Uint256
}

func(this *Key) ToBytes() []byte{
	sink := &common.ZeroCopySink{}
	this.Serialize(sink)
	return sink.Bytes()
}

func(this *Key) Serialize(sink *common.ZeroCopySink){
	sink.WriteByte(byte(this.Type))
	sink.WriteHash(this.Hash)
}