package block

import "github.com/ontio/ontology/common"

type Key struct{
	Type Type
	Sign common.Uint256
}

func(this *Key) GetType() Type{
	return this.Type
}

func (this *Key) SerializeKey()  []byte{
	sink := &common.ZeroCopySink{}
	this.Serialize(sink)
	return sink.Bytes()
}

func(this *Key) Serialize(sink *common.ZeroCopySink){
	sink.WriteByte(byte(this.Type))
	sink.WriteHash(this.Sign)
}
