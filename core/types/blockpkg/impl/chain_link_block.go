package impl

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/ontio/ontology/common"
)

type ChainLinkBlock struct {
	Block
	previousHash interfaces.IKey
}

func (this *ChainLinkBlock) GetOrigin() interfaces.IKey {
	return this.body.GetOrigin()
}

func (this *ChainLinkBlock) ValueToBytes() []byte {
	sink := &common.ZeroCopySink{}
	sink.WriteVarBytes(this.header.GetSign())
	this.body.SerializeData(sink)
	this.previousHash.Serialize(sink)
	return sink.Bytes()
}

func (this *ChainLinkBlock) GetPreviousHash(prevHashType blockpkg.PrevHashType) interfaces.IKey {
	return this.previousHash
}

func (this *ChainLinkBlock) SetPreviousHash(prevHashType blockpkg.PrevHashType, key interfaces.IKey) {
	this.previousHash = key
}
