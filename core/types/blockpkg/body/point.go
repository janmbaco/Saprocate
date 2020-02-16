package body

import (
	"crypto/sha256"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/ontio/ontology/common"
)

type Point struct {
	Origin    *header.Key
	To        *header.Key
	Timestamp uint64
	Nonce     uint32
	Sign      []byte
	Transfers []*Transfer
}

func (this *Point) SerializeData(sink *common.ZeroCopySink) {
	this.Origin.Serialize(sink)
	this.To.Serialize(sink)
	sink.WriteUint64(this.Timestamp)
	sink.WriteUint32(this.Nonce)
	sink.WriteVarBytes(this.Sign)
}

func (this *Point) GetDataSigned() []byte {
	sink := &common.ZeroCopySink{}
	this.Origin.Serialize(sink)
	this.To.Serialize(sink)
	sink.WriteUint64(this.Timestamp)
	sink.WriteUint32(this.Nonce)
	return sink.Bytes()
}
func (this *Point) GetOrigin() *header.Key {
	return this.Origin
}

func (this *Point) GetDestiny() *header.Key {
	return this.To
}

func (this *Point) GetSign() []byte {
	return this.Sign
}

func (this *Point) GetHash() common.Uint256 {
	signSum := sha256.Sum256(this.Sign)
	ui256, _ := common.Uint256ParseFromBytes(signSum[:])
	return ui256
}
