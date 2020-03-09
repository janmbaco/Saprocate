package impl

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/ontio/ontology/common"
)

type DoubleChainLinkBlock struct {
	ChainLinkBlock
	secondPreviousHash interfaces.IKey
}

func (this *DoubleChainLinkBlock) ValueToBytes() []byte {
	sink := &common.ZeroCopySink{}
	sink.WriteVarBytes(this.header.GetSign())
	this.body.SerializeData(sink)
	this.previousHash.Serialize(sink)
	this.secondPreviousHash.Serialize(sink)
	return sink.Bytes()
}

func (this *DoubleChainLinkBlock) GetPreviousHash(prevHashType blockpkg.PrevHashType) interfaces.IKey {
	prevHash := this.previousHash
	if prevHashType == blockpkg.SecondPrvHash {
		prevHash = this.secondPreviousHash
	}
	return prevHash
}

func (this *DoubleChainLinkBlock) SetPreviousHash(prevHashType blockpkg.PrevHashType, key interfaces.IKey) {
	if prevHashType == blockpkg.SecondPrvHash {
		this.secondPreviousHash = key
	} else {
		this.previousHash = key
	}
}
