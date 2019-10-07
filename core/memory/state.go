package memory

import "github.com/ontio/ontology/common"

type State struct {
}

func(this *State) GetLastHash() *common.Uint256 {
	return nil
}

func(this *State) GetOwner(blockHash *common.Uint256) *common.Uint256{
	return nil
}
