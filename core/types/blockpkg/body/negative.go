package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/ontio/ontology/common"
)

type Negative struct {
	positiveBlock interfaces.IBlock
}

func NewNegative(positiveBlock interfaces.IBlock) *Negative {
	return &Negative{positiveBlock: positiveBlock}
}

func (this *Negative) SerializeData(sink *common.ZeroCopySink) {
	sink.WriteVarBytes(this.positiveBlock.ValueToBytes())
}

func (this *Negative) GetOrigin() interfaces.IKey {
	return this.positiveBlock.GetOrigin()
}

func (this *Negative) GetPositiveBlock() interfaces.IBlock {
	return this.positiveBlock
}
