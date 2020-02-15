package blockpkg

import "github.com/ontio/ontology/common"

type ChainLinkBlock struct {
	Block
	PrevHashKey *Key
}

func (this *ChainLinkBlock) GetOrigin() *Key {
	return this.Body.GetOrigin()
}

func (this *ChainLinkBlock) ValueToBytes() []byte {
	sink := &common.ZeroCopySink{}
	sink.WriteVarBytes(this.Header.Sign)
	this.Body.SerializeData(sink)
	this.PrevHashKey.Serialize(sink)
	return sink.Bytes()
}
