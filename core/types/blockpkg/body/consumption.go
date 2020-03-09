package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/ontio/ontology/common"
)

type Consumption struct {
	positives []interfaces.IBlock
}

func NewConsumption(positives []interfaces.IBlock) *Consumption {
	return &Consumption{positives: positives}
}

func (this *Consumption) SerializeData(sink *common.ZeroCopySink) {
	sink.WriteVarUint(uint64(len(this.positives)))
	for _, positive := range this.positives {
		sink.WriteVarBytes(positive.ValueToBytes())
	}
}

func (this *Consumption) GetOrigin() interfaces.IKey {
	return this.positives[0].GetOrigin()
}

func (this *Consumption) GetPositiveBlocks() []interfaces.IBlock {
	return this.positives
}
