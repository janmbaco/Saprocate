package types

import (
	"crypto/rsa"
	"encoding/asn1"
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/ontio/ontology/common"
)

type OriginBlock struct{
	block.Key
	PubilcKey rsa.PublicKey
}

func(origin *OriginBlock) SerializeValue() []byte{
	sink := common.ZeroCopySink{}
	pubBytes, err := asn1.Marshal(origin.PubilcKey)
	cross.TryPanic(err)
	sink.WriteBytes(pubBytes)
	return sink.Bytes()
}
