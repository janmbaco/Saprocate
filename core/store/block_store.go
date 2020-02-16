package store

import (
	"crypto/rsa"
	"encoding/asn1"
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/body"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/impl"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/ontio/ontology/common"
	"io"
)

func BlockFromBytes(key *header.Key, value []byte) blockpkg.Interface {
	var result blockpkg.Interface
	switch key.Type {
	case blockpkg.Origin:
		result = newOriginBlock(key, value)
	case blockpkg.Positive:
		result = newPositiveBlock(key, value)
	case blockpkg.Negative:
		result = newNegativeBlock(key, value)
	case blockpkg.Transfer:
		result = newTransferBlock(key, value)
	case blockpkg.Pay:
		result = newPayBlock(key, value)
	}
	return result
}

func KeyFromBytes(raw []byte) *header.Key {
	return getSourceKey(common.NewZeroCopySource(raw))
}

func newOriginBlock(key *header.Key, value []byte) *impl.Block {
	source := common.NewZeroCopySource(value)
	header := getSourceHeader(key, source)
	buff, _, _, eof := source.NextVarBytes()
	tryEof(eof)
	pk := new(rsa.PublicKey)
	_, err := asn1.Unmarshal(buff, pk)
	cross.TryPanic(err)
	body := &body.Origin{
		PublicKey: pk,
	}
	return &impl.Block{
		Header: header,
		Body:   body,
	}
}

func newPositiveBlock(key *header.Key, value []byte) *impl.ChainLinkBlock {
	source := common.NewZeroCopySource(value)
	header := getSourceHeader(key, source)
	point := getSourcePoint(source)
	body :=  &body.Positive{
		Point:          point,
	}
	prev := getSourceKey(source)
	return &impl.ChainLinkBlock{
		Block:       impl.Block{
			Header: header,
			Body:   body,
		},
		PrevHashKey: prev,
	}
}

func newNegativeBlock(key *header.Key, value []byte) *impl.ChainLinkBlock {
	source := common.NewZeroCopySource(value)
	header := getSourceHeader(key, source)
	positiveBlock := getSourceKey(source)
	body :=  &body.Negative{
		PositiveBlockKey: positiveBlock,
	}
	prev := getSourceKey(source)
	return &impl.ChainLinkBlock{
		Block:       impl.Block{
			Header: header,
			Body:   body,
		},
		PrevHashKey: prev,
	}
}

func newTransferBlock(key *header.Key, value []byte) *impl.ChainLinkBlock {
	source := common.NewZeroCopySource(value)
	header := getSourceHeader(key, source)
	from := getSourceKey(source)
	to := getSourceKey(source)
	var coins []*body.Point
	m, eof := source.NextUint64()
	tryEof(eof)
	for i := 0; i< int(m); i++{
		coins = append(coins, getSourcePoint(source))
	}
	body := &body.Transfer{
		From:           from,
		To:             to,
		Points:         coins,
	}
	prev := getSourceKey(source)
	return &impl.ChainLinkBlock{
		Block:       impl.Block{
			Header: header,
			Body:   body,
		},
		PrevHashKey: prev,
	}
}

func newPayBlock(key *header.Key, value []byte) *impl.ChainLinkBlock {
	source := common.NewZeroCopySource(value)
	header := getSourceHeader(key, source)
	from := getSourceKey(source)
	var points []*body.Point
	m, eof := source.NextUint64()
	tryEof(eof)
	for i := 0; i< int(m); i++{
		points = append(points, getSourcePoint(source))
	}
	body := &body.Pay{
		From:           from,
		Points:         points,
	}
	prev := getSourceKey(source)
	return &impl.ChainLinkBlock{
		Block:       impl.Block{
			Header: header,
			Body:   body,
		},
		PrevHashKey: prev,
	}
}

func getSourceKey(source *common.ZeroCopySource) *header.Key {
	t, eof := source.NextByte()
	tryEof(eof)
	hash, eof := source.NextHash()
	tryEof(eof)
	return &header.Key{
		Type: blockpkg.Type(t),
		Hash: hash,
	}
}

func getSourceHeader(key *header.Key, source *common.ZeroCopySource) *header.Header {
	buff, _, _, eof := source.NextVarBytes()
	tryEof(eof)
	return &header.Header{
		Key:       key,
		Sign:      buff,
	}
}


func getSourcePoint(source *common.ZeroCopySource) *body.Point {
	origin := getSourceKey(source)
	to := getSourceKey(source)
	timeStamp, eof := source.NextUint64()
	tryEof(eof)
	sign, _, _, eof := source.NextVarBytes()
	tryEof(eof)
	return &body.Point{
		Origin: origin,
		To: to,
		Timestamp: timeStamp,
		Sign: sign,
	}
}

func tryEof(eof bool) {
	if eof {
		cross.TryPanic(io.ErrUnexpectedEOF)
	}
}
