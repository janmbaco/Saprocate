package types

import (
	"github.com/ontio/ontology/common"
)

type Header struct {
	PrevHash common.Uint256
	Timestamp uint64 //in nanoseconds
	Nonce uint64
	SignOrigin common.Uint256 //hash publick key
}

func (this *Header) Serialization(sink *common.ZeroCopySink) {
	sink.WriteHash(this.PrevHash)
	sink.WriteUint64(this.Timestamp)
	sink.WriteUint64(this.Nonce)
	sink.WriteHash(this.SignOrigin)
}

func (this *Header) Deserialization (source *common.ZeroCopySource)  {
	var eof bool
	this.PrevHash, eof = source.NextHash()
	tryEof(eof)
	this.Timestamp, eof = source.NextUint64()
	tryEof(eof)
	this.Nonce, eof = source.NextUint64()
	tryEof(eof)
	this.SignOrigin, eof = source.NextHash()
	tryEof(eof)
}

