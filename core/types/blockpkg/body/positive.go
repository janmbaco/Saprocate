package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/ontio/ontology/common"
)

type Positive struct{
	Point    *Point
}

func(this *Positive) SerializeData(sink *common.ZeroCopySink) {
	this.Point.Serilize(sink)
}

func(this *Positive) GetOrigin() *header.Key {
	return this.Point.To
}

func(this *Positive) GetPoints() []*Point {
	return []*Point{this.Point}
}
