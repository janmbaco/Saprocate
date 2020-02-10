package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/ontio/ontology/common"
)

type Negative struct{
	PositiveBlockKey *blockpkg.Key
}

func(this *Negative) SerializeData(sink *common.ZeroCopySink) {
	this.PositiveBlockKey.Serialize(sink)
}

func(this *Negative) GetOrigin() *blockpkg.Key{
	return this.PositiveBlockKey
}

