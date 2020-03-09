package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/ontio/ontology/common"
)

type Positive struct {
	point interfaces.IPoint
	to    interfaces.IKey
}

func NewPositive(point interfaces.IPoint, to interfaces.IKey) *Positive {
	return &Positive{point: point, to: to}
}

func (this *Positive) SerializeData(sink *common.ZeroCopySink) {
	this.point.Serialize(sink)
	this.to.Serialize(sink)
}

func (this *Positive) GetOrigin() interfaces.IKey {
	return this.to
}

func (this *Positive) GetPoint() interfaces.IPoint {
	return this.point
}
