package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/ontio/ontology/common"
)

type Pay struct{
	From     *header.Key
	Points   []*Point
}

func(this *Pay) SerializeData(sink *common.ZeroCopySink) {
	this.From.Serialize(sink)
	sink.WriteVarUint(uint64(len(this.Points)))
	for _, coin := range this.Points {
		coin.Serilize(sink)
	}
}

func(this *Pay) GetOrigin() *header.Key {
	return this.From
}

func(this *Pay) GetPoints() []*Point {
	return this.Points
}
