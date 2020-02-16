package types

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
)

type (
	Summary struct {
		Owner *header.Key
		PointsCards map[header.Key]uint
	})
