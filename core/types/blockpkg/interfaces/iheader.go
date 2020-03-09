package interfaces

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/ontio/ontology/common"
)

type (
	IKey interface {
		GetHash() common.Uint256
		GetType() blockpkg.BlockType
		Serialize(*common.ZeroCopySink)
		ToBytes() []byte
	}
	IHeader interface {
		IKey
		GetSign() []byte
		SetSign([]byte)
		GetKey() IKey
	}
)
