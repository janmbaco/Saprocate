package memory

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/Point"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/ontio/ontology/common"
)

type Memory struct {
	Point map[common.Uint256]State
}

func NewMemory() *Memory {
	return &Memory{Point: make(map[common.Uint256]State)}
}

func(this *Memory) SetPoint(t header.Type, point Point.Interface){
	hash := point.GetHash()
	switch t {
	case header.Numu:
		this.Point[hash] = Stored
	case header.Transfer:
		this.Point[hash] = Transferred
	case header.Pay:
		this.Point[hash] = Paid
	}
}


