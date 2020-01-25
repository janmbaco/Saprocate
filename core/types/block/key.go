package block

import "github.com/ontio/ontology/common"

type Key struct{
	Type Type
	Sign common.Uint256
}

func (key *Key) SerializeKey()  []byte{
	sink := common.ZeroCopySink{}
	sink.WriteByte(byte(key.Type))
	sink.WriteHash(key.Sign)
	return sink.Bytes()
}
