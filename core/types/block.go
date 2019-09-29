package types

import (
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/ontio/ontology/common"
)

type Block struct{
	Raw *[]byte //block saved
	Hash common.Uint256
	Header *Header
	Body	Body
}

func(this *Block) ToArray() []byte{
	if this.Raw == nil{
		sink := &common.ZeroCopySink{}
		sink.WriteHash(this.Hash)
		this.Header.Serialization(sink)
		this.Body.Serialization(sink)
		bytes := sink.Bytes()
		this.Raw = &bytes
	}
	return *this.Raw
}

func(this *Block) getFromRaw(){
	if this.Raw == nil {
		cross.TryPanic("Nil blok")
	}

	source:= common.NewZeroCopySource(*this.Raw)
	var eof bool
	this.Hash, eof = source.NextHash()
	tryEof(eof)
	this.Header.Deserialization(source)
	this.Body.Deserialization(source)

}

func(this *Block) Type() BlockType {
	if this.Body == nil && this.Raw == nil {
		cross.TryPanic("This is a nil Block")
	}
	if this.Body == nil  {
		this.getFromRaw()
	}
	return this.Body.GetType()
}

func(this *Block) GetDataToSign() []byte{
	sink := &common.ZeroCopySink{}
	this.Header.Serialization(sink)
	this.Body.Serialization(sink)
	return sink.Bytes()
}

