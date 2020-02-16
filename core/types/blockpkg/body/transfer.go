package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/impl"
	"github.com/ontio/ontology/common"
)

type Transfer struct{
	impl.ChainLinkBlock
	From     *header.Key
	To       *header.Key
	Points   []*Point
}

func(this *Transfer) SerializeData(sink *common.ZeroCopySink) {
	this.From.Serialize(sink)
	this.To.Serialize(sink)
	sink.WriteVarUint(uint64(len(this.Points)))
	for  _ , point := range this.Points {
		point.Serilize(sink)
	}
}

func(this *Transfer) GetOrigin() *header.Key {
	return this.From
}

func(this *Transfer) GetPoints() []*Point {
	return this.Points
}


