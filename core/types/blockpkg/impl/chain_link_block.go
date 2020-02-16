package impl

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/ontio/ontology/common"
)

type ChainLinkBlock struct {
	Block
	PrevHashKey *header.Key
}

func (this *ChainLinkBlock) GetOrigin() *header.Key {
	return this.Body.GetOrigin()
}

func (this *ChainLinkBlock) ValueToBytes() []byte {
	sink := &common.ZeroCopySink{}
	sink.WriteVarBytes(this.Header.Sign)
	this.Body.SerializeData(sink)
	this.PrevHashKey.Serialize(sink)
	return sink.Bytes()
}

func(this *ChainLinkBlock) GetDataSigned() []byte{
	sink := &common.ZeroCopySink{}
	this.Body.SerializeData(sink)
	this.PrevHashKey.Serialize(sink)
	return sink.Bytes()
}
