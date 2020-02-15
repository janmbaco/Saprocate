package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/ontio/ontology/common"
)

type Pay struct{
	From     *blockpkg.Key
	Points   []*blockpkg.Point
}

func(this *Pay) SerializeData(sink *common.ZeroCopySink) {
	this.From.Serialize(sink)
	sink.WriteVarUint(uint64(len(this.Points)))
	for _, coin := range this.Points {
		coin.Serilize(sink)
	}
}

func(this *Pay) GetOrigin() *blockpkg.Key{
	return this.From
}

func(this *Pay) GetPoints() []*blockpkg.Point{
	return this.Points
}
