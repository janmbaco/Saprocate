package body

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/ontio/ontology/common"
)

type Voucher struct {
	point       interfaces.IPoint
	negativeKey interfaces.IKey
	positiveKey interfaces.IKey
	origin      interfaces.IKey
}

func NewVoucher(point interfaces.IPoint, negativeKey interfaces.IKey, positiveKey interfaces.IKey, origin interfaces.IKey) *Voucher {
	persistentHashNegative := header.NewKey(blockpkg.Negative, common.UINT256_EMPTY)
	if negativeKey != nil {
		persistentHashNegative = negativeKey.(*header.Key)
	}
	return &Voucher{point: point, negativeKey: persistentHashNegative, positiveKey: positiveKey, origin: origin}
}

func (this *Voucher) SerializeData(sink *common.ZeroCopySink) {
	this.point.Serialize(sink)
	this.negativeKey.Serialize(sink)
	this.positiveKey.Serialize(sink)
	this.origin.Serialize(sink)
}

func (this *Voucher) GetOrigin() interfaces.IKey {
	return this.origin
}

func (this *Voucher) GetPoint() interfaces.IPoint {
	return this.point
}

func (this *Voucher) HasNegative() bool {
	return this.negativeKey.GetHash() != common.UINT256_EMPTY
}

func (this *Voucher) GetNegativeKey() interfaces.IKey {
	return this.negativeKey
}

func (this *Voucher) GetPositivekey() interfaces.IKey {
	return this.positiveKey
}
