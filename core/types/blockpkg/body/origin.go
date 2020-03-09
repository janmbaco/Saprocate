package body

import (
	"crypto/rsa"
	"encoding/asn1"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/ontio/ontology/common"
)

type Origin struct {
	publicKey *rsa.PublicKey
}

func NewOrigin(publicKey *rsa.PublicKey) *Origin {
	return &Origin{publicKey: publicKey}
}

func (this *Origin) SerializeData(sink *common.ZeroCopySink) {
	pubBytes, err := asn1.Marshal(*this.publicKey)
	cross.TryPanic(err)
	sink.WriteVarBytes(pubBytes)
}

func (this *Origin) GetPublicKey() *rsa.PublicKey {
	return this.publicKey
}

func (this *Origin) GetOrigin() interfaces.IKey {
	return nil
}
