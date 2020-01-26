package block

import "github.com/ontio/ontology/common"

type Coin struct{
	Timestamp uint64
	Origin common.Uint256
	Sign common.Uint256
}

func(coin *Coin) Serilize(sink *common.ZeroCopySink){
	sink.WriteUint64(coin.Timestamp)
	sink.WriteHash(coin.Origin)
	sink.WriteHash(coin.Sign)
}
