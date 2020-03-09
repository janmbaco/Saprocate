package interfaces

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
)

type IBlock interface {
	GetHeader() IHeader
	GetOrigin() IKey
	GetSign() []byte
	SetSign([]byte)
	GetBody() IBody
	GetDataSigned() []byte
	KeyToBytes() []byte
	ValueToBytes() []byte
	GetPreviousHash(blockpkg.PrevHashType) IKey
	SetPreviousHash(blockpkg.PrevHashType, IKey)
}
