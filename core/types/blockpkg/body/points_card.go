package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/ontio/ontology/common"
)

type PointCard struct{
	From     *header.Key
	Points   []*Point
}

func(this *PointCard) SerializeData(sink *common.ZeroCopySink) {
	this.From.Serialize(sink)
	sink.WriteVarUint(uint64(len(this.Points)))
	for _, coin := range this.Points {
		coin.SerializeData(sink)
	}
}

func(this *PointCard) GetOrigin() *header.Key {
	return this.From
}

func(this *PointCard) GetPoints() []*Point {
	return this.Points
}
