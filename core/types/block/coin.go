package block

import "github.com/ontio/ontology/common"

type Coin struct{
	Origin *Key
	Timestamp uint64
	Sign common.Uint256
}

func(this *Coin) Serilize(sink *common.ZeroCopySink){
	this.Origin.Serialize(sink)
	sink.WriteUint64(this.Timestamp)
	sink.WriteHash(this.Sign)
}
