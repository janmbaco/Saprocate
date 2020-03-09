package types

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
)

type Summary struct {
	Owner       interfaces.IKey
	PointsCards map[interfaces.IKey]uint
}

func NewSummary(owner interfaces.IKey) *Summary {
	return &Summary{
		Owner:       owner,
		PointsCards: make(map[interfaces.IKey]uint),
	}
}
