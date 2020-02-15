package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/ontio/ontology/common"
)

type Positive struct{
	Point    *blockpkg.Point
}

func(this *Positive) SerializeData(sink *common.ZeroCopySink) {
	this.Point.Serilize(sink)
}

func(this *Positive) GetOrigin() *blockpkg.Key{
	return this.Point.To
}

func(this *Positive) GetPoints() []*blockpkg.Point{
	return []*blockpkg.Point{this.Point}
}
