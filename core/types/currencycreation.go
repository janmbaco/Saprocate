package types

import (
	"crypto/rsa"
	"encoding/asn1"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/ontio/ontology/common"
)

//block type for creation of currency
type CurrencyCreation struct {
	PublicKey *rsa.PublicKey // key to verify the coin creation
}

func (this *CurrencyCreation) Serialization(sink *common.ZeroCopySink) {
	pubBytes, err := asn1.Marshal(this.PublicKey)
	cross.TryPanic(err)
	sink.WriteBytes(pubBytes)
}

func (this *CurrencyCreation) Deserialization (source *common.ZeroCopySource)  {
	buf, _, _, eof :=  source.NextVarBytes()
	tryEof(eof)
	_, err := asn1.Unmarshal(buf, this.PublicKey)
	cross.TryPanic(err)
}

func(this *CurrencyCreation) GetType() BlockType{
	return CurrencyCreator
}