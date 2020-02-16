package blockpkg

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg/body"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/ontio/ontology/common"
)

type (
	Interface interface{
		GetType() Type
		GetOrigin() *header.Key
		GetSign() []byte
		GetDataSigned() []byte
		KeyToBytes() []byte
		ValueToBytes() []byte
	}
	)
