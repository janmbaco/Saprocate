package types

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
)

type (
	Summary struct {
		Owner *blockpkg.Key
		PointsCards map[blockpkg.Key]uint
	})
