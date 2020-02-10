package blockpkg

import "github.com/ontio/ontology/common"

type (
	Interface interface{
		GetType() Type
		GetOrigin() *Key
		GetSign() []byte
		GetDataSigned() []byte
		KeyToBytes() []byte
		ValueToBytes() []byte
	}
	BodyBlock interface {
		SerializeData(*common.ZeroCopySink)
		GetOrigin() *Key
	}
	PointsBody interface{
		GetPoints() []*Point
	})
