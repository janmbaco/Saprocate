package impl

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/body"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/ontio/ontology/common"
)

type Block struct{
	Header *header.Header
	Body   body.Interface
}

func(this *Block) GetType() header.Type {
	return this.Header.Key.Type
}

func(this *Block) GetOrigin() *header.Key {
	return this.Header.Key
}

func(this *Block) GetSign() []byte{
	return this.Header.Sign
}

func(this *Block) GetDataSigned() []byte{
	sink := &common.ZeroCopySink{}
	this.Body.SerializeData(sink)
	return sink.Bytes()
}

func(this *Block) KeyToBytes() []byte{
	return this.Header.Key.ToBytes()
}

func(this *Block) ValueToBytes() []byte{
	sink := &common.ZeroCopySink{}
	sink.WriteVarBytes(this.Header.Sign)
	this.Body.SerializeData(sink)
	return sink.Bytes()
}