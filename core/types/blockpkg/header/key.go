package header

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/ontio/ontology/common"
)

type Key struct {
	blockType blockpkg.BlockType
	hash      common.Uint256
}

func NewKey(typ blockpkg.BlockType, hash common.Uint256) *Key {
	return &Key{blockType: typ, hash: hash}
}

func (this *Key) GetType() blockpkg.BlockType {
	return this.blockType
}

func (this *Key) GetHash() common.Uint256 {
	return this.hash
}

func (this *Key) ToBytes() []byte {
	sink := &common.ZeroCopySink{}
	this.Serialize(sink)
	return sink.Bytes()
}

func (this *Key) Serialize(sink *common.ZeroCopySink) {
	sink.WriteByte(byte(this.blockType))
	sink.WriteHash(this.hash)
}
