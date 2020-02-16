package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/ontio/ontology/common"
)

type Negative struct{
	PositiveBlockKey *header.Key
}

func(this *Negative) SerializeData(sink *common.ZeroCopySink) {
	this.PositiveBlockKey.Serialize(sink)
}

func(this *Negative) GetOrigin() *header.Key {
	return this.PositiveBlockKey
}

