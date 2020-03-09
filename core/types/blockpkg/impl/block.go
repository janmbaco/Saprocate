package impl

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/ontio/ontology/common"
)

type Block struct {
	header interfaces.IHeader
	body   interfaces.IBody
}

func (this *Block) GetType() blockpkg.BlockType {
	return this.header.GetType()
}

func (this *Block) GetHeader() interfaces.IHeader {
	return this.header
}

func (this *Block) GetOrigin() interfaces.IKey {
	return this.header.GetKey()
}

func (this *Block) GetBody() interfaces.IBody {
	return this.body
}

func (this *Block) GetSign() []byte {
	return this.header.GetSign()
}

func (this *Block) SetSign(sign []byte) {
	this.header.SetSign(sign)
}

func (this *Block) GetDataSigned() []byte {
	sink := &common.ZeroCopySink{}
	this.body.SerializeData(sink)
	return sink.Bytes()
}

func (this *Block) KeyToBytes() []byte {
	return this.header.GetKey().ToBytes()
}

func (this *Block) ValueToBytes() []byte {
	sink := &common.ZeroCopySink{}
	sink.WriteVarBytes(this.header.GetSign())
	this.body.SerializeData(sink)
	return sink.Bytes()
}

func (this *Block) GetPreviousHash(prevHashType blockpkg.PrevHashType) interfaces.IKey {
	return nil
}

func (this *Block) SetPreviousHash(prevHashType blockpkg.PrevHashType, key interfaces.IKey) {
	//do nothing
}
