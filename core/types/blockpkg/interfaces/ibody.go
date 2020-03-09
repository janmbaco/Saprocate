package interfaces

import (
	"github.com/ontio/ontology/common"
)

type IBody interface {
	SerializeData(*common.ZeroCopySink)
	GetOrigin() IKey
}
