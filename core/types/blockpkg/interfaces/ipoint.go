package interfaces

import "github.com/ontio/ontology/common"

type IPoint interface {
	Serialize(sink *common.ZeroCopySink)
	GetDataSigned() []byte
	GetOrigin() IKey
	GetSign() []byte
	SetSign(sign []byte)
	GetHash() common.Uint256
}
