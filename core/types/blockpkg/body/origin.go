package body

import (
	"crypto/rsa"
	"encoding/asn1"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/ontio/ontology/common"
)

type Origin struct{
	PublicKey *rsa.PublicKey
}

func(this *Origin) SerializeData(sink *common.ZeroCopySink){
	pubBytes, err := asn1.Marshal(*this.PublicKey)
	cross.TryPanic(err)
	sink.WriteVarBytes(pubBytes)
}

func(this *Origin) GetOrigin() *header.Key {
	return nil
}


