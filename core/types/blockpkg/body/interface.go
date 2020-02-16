package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/ontio/ontology/common"
)

type (
	Interface interface {
		SerializeData(*common.ZeroCopySink)
		GetOrigin() *header.Key
	}
	PointsBody interface {
		GetPoints() []*Point
	}
)
