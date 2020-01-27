package types

import (
	"crypto/rsa"
	"encoding/asn1"
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/ontio/ontology/common"
	"io"
)

func NewFromBytes(key *block.Key, value []byte) block.Interface {
	return getBlock(key, value)
}

func KeyFromBytes(raw []byte) *block.Key {
	return getSourceKey(common.NewZeroCopySource(raw))
}

func getBlock(key *block.Key, value []byte) block.Interface {
	var result block.Interface
	switch key.Type {
	case block.Origin:
		result = newOriginBlock(key, value)
	case block.Positive:
		result = newPositiveBlock(key, value)
	case block.Negative:
		result = newNegativeBlock(key, value)
	case block.Transfer:
		result =newTransferBlock(key, value)
	case block.Pay:
		result = newPayBlock(key, value)
	}
	return result
}

func newOriginBlock(key *block.Key, value []byte) *OriginBlock {
	source := common.NewZeroCopySource(value)
	buff, _, _, eof := source.NextVarBytes()
	tryEof(eof)
	var pk *rsa.PublicKey
	_, err := asn1.Unmarshal(buff, pk)
	cross.TryPanic(err)
	return &OriginBlock{
		Key: *key,
		PubilcKey: pk,
	}
}

func newPositiveBlock(key *block.Key, value []byte) *PositiveBlock {
	source := common.NewZeroCopySource(value)
	prev := getSourceKey(source)
	coin := getSourceCoin(source)
	return &PositiveBlock{
		Key: *key,
		Previous: prev,
		Coin: coin,
	}
}

func newNegativeBlock(key *block.Key, value []byte) *NegativeBlock {
	source := common.NewZeroCopySource(value)
	prev := getSourceKey(source)
	positiveBlock := getSourceKey(source)
	return &NegativeBlock{
		Key: *key,
		Previous:prev,
		PositiveBlock:positiveBlock,
	}
}

func newTransferBlock(key *block.Key, value []byte) *TransferBlock {
	source := common.NewZeroCopySource(value)
	prev := getSourceKey(source)
	from := getSourceKey(source)
	to := getSourceKey(source)
	var coins []*block.Coin
	m, eof := source.NextUint64()
	tryEof(eof)
	for i := 0; i< int(m); i++{
		coins = append(coins, getSourceCoin(source))
	}
	return &TransferBlock{
		Key: *key,
		Previous:prev,
		From:from,
		To:to,
		Coins: coins,
	}
}

func newPayBlock(key *block.Key, value []byte) *PayBlock {
	source := common.NewZeroCopySource(value)
	prev := getSourceKey(source)
	from := getSourceKey(source)
	var coins []*block.Coin
	m, eof := source.NextUint64()
	tryEof(eof)
	for i := 0; i< int(m); i++{
		coins = append(coins, getSourceCoin(source))
	}
	return &PayBlock{
		Key: *key,
		Previous:prev,
		From:from,
		Coins: coins,
	}
}

func getSourceKey(source *common.ZeroCopySource) *block.Key {
	t, eof := source.NextByte()
	tryEof(eof)
	hash, eof := source.NextHash()
	tryEof(eof)
	return &block.Key{
		Type: block.Type(t),
		Sign: hash,
	}
}

func getSourceCoin(source *common.ZeroCopySource) *block.Coin {
	origin := getSourceKey(source)
	timeStamp, eof := source.NextUint64()
	tryEof(eof)
	sign, eof := source.NextHash()
	tryEof(eof)
	return &block.Coin{
		Origin: origin,
		Timestamp: timeStamp,
		Sign: sign,
	}
}

func tryEof(eof bool) {
	if eof {
		cross.TryPanic(io.ErrUnexpectedEOF)
	}
}
