package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/ontio/ontology/common"
)

type Point struct{
	Origin *header.Key
	To *header.Key
	Timestamp uint64
	Sign []byte
}

func(this *Point) Serilize(sink *common.ZeroCopySink){
	this.Origin.Serialize(sink)
	this.To.Serialize(sink)
	sink.WriteUint64(this.Timestamp)
	sink.WriteVarBytes(this.Sign)
}

func(this *Point) GetDataSigned() []byte{
	sink := &common.ZeroCopySink{}
	this.Origin.Serialize(sink)
	this.To.Serialize(sink)
	sink.WriteUint64(this.Timestamp)
	return sink.Bytes()
}
