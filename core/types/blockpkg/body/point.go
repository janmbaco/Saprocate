package body

import (
	"crypto/sha256"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/ontio/ontology/common"
)

type Point struct {
	origin     interfaces.IKey
	timestamp  uint64
	nonce      uint32
	expireDate uint64
	sign       []byte
}

func NewPoint(origin interfaces.IKey, timestamp uint64, nonce uint32, expireDate uint64) *Point {
	return &Point{origin: origin, timestamp: timestamp, nonce: nonce, expireDate: expireDate}
}

func (this *Point) Serialize(sink *common.ZeroCopySink) {
	this.origin.Serialize(sink)
	sink.WriteUint64(this.timestamp)
	sink.WriteUint32(this.nonce)
	sink.WriteUint64(this.expireDate)
	sink.WriteVarBytes(this.sign)
}

func (this *Point) GetDataSigned() []byte {
	sink := &common.ZeroCopySink{}
	this.origin.Serialize(sink)
	sink.WriteUint64(this.timestamp)
	sink.WriteUint32(this.nonce)
	sink.WriteUint64(this.expireDate)
	return sink.Bytes()
}

func (this *Point) GetOrigin() interfaces.IKey {
	return this.origin
}

func (this *Point) GetSign() []byte {
	return this.sign
}

func (this *Point) SetSign(sign []byte) {
	this.sign = sign
}

func (this *Point) GetHash() common.Uint256 {
	signSum := sha256.Sum256(this.sign)
	ui256, _ := common.Uint256ParseFromBytes(signSum[:])
	return ui256
}
