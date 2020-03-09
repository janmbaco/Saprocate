package header

import (
	"crypto/sha256"
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/ontio/ontology/common"
)

type Header struct {
	Key
	sign []byte
}

func NewHeader(blockType blockpkg.BlockType, sign []byte) *Header {
	header := &Header{}
	header.blockType = blockType
	header.sign = sign
	return header
}

func (this *Header) GetSign() []byte {
	return this.sign
}

func (this *Header) SetSign(sign []byte) {
	this.sign = sign
}

func (this *Header) GetHash() common.Uint256 {
	signSum := sha256.Sum256(this.sign)
	ui256, _ := common.Uint256ParseFromBytes(signSum[:])
	return ui256
}

func (this *Header) GetType() blockpkg.BlockType {
	return this.blockType
}

func (this *Header) GetKey() interfaces.IKey {
	return NewKey(this.blockType, this.GetHash())
}

func (this *Header) ToBytes() []byte {
	sink := &common.ZeroCopySink{}
	this.Serialize(sink)
	return sink.Bytes()
}

func (this *Header) Serialize(sink *common.ZeroCopySink) {
	sink.WriteByte(byte(this.blockType))
	sink.WriteHash(this.GetHash())
}
