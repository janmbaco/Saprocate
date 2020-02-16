package blockpkg

import "github.com/janmbaco/Saprocate/core/types/blockpkg/header"

type Interface interface {
	GetType() header.Type
	GetOrigin() *header.Key
	GetSign() []byte
	GetDataSigned() []byte
	KeyToBytes() []byte
	ValueToBytes() []byte
}
