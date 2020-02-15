package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/ontio/ontology/common"
)

type Transfer struct{
	blockpkg.ChainLinkBlock
	From     *blockpkg.Key
	To       *blockpkg.Key
	Points   []*blockpkg.Point
}

func(this *Transfer) SerializeData(sink *common.ZeroCopySink) {
	this.From.Serialize(sink)
	this.To.Serialize(sink)
	sink.WriteVarUint(uint64(len(this.Points)))
	for  _ , point := range this.Points {
		point.Serilize(sink)
	}
}

func(this *Transfer) GetOrigin() *blockpkg.Key{
	return this.From
}

func(this *Transfer) GetPoints() []*blockpkg.Point{
	return this.Points
}


