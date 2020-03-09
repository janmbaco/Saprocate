package store

import (
	"crypto/rsa"
	"encoding/asn1"
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/body"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/impl"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/ontio/ontology/common"
	"io"
)

func BlockFromBytes(key interfaces.IKey, value []byte) interfaces.IBlock {
	var result interfaces.IBlock
	switch key.GetType() {
	case blockpkg.Origin:
		result = newOriginBlock(value)
	case blockpkg.Positive:
		result = newPositiveBlock(value)
	case blockpkg.Negative:
		result = newNegativeBlock(value)
	case blockpkg.Voucher:
		result = newVoucherBlock(value)
	case blockpkg.Consumption:
		result = newConsumptionBlock(value)
	}
	return result
}

func KeyFromBytes(raw []byte) interfaces.IKey {
	return getSourceKey(common.NewZeroCopySource(raw))
}

func newOriginBlock(value []byte) interfaces.IBlock {
	source := common.NewZeroCopySource(value)
	sign := getSourceSign(source)
	buff, _, _, eof := source.NextVarBytes()
	tryEof(eof)
	pk := new(rsa.PublicKey)
	_, err := asn1.Unmarshal(buff, pk)
	cross.TryPanic(err)
	originBlock := impl.NewOriginBlock(pk)
	originBlock.SetSign(sign)
	return originBlock
}

func newPositiveBlock(value []byte) interfaces.IBlock {
	source := common.NewZeroCopySource(value)
	return getSourcePositiveBlock(source)
}

func newNegativeBlock(value []byte) interfaces.IBlock {
	source := common.NewZeroCopySource(value)
	sign := getSourceSign(source)
	positiveBlock := getSourcePositiveBlock(source)
	firstPrev := getSourceKey(source)
	secondPrev := getSourceKey(source)
	negativeBlock := impl.NewNegativeblock(positiveBlock)
	negativeBlock.SetSign(sign)
	negativeBlock.SetPreviousHash(blockpkg.FirstPrevHash, firstPrev)
	negativeBlock.SetPreviousHash(blockpkg.SecondPrvHash, secondPrev)
	return negativeBlock
}

func newVoucherBlock(value []byte) interfaces.IBlock {
	source := common.NewZeroCopySource(value)
	sign := getSourceSign(source)
	point := getSourcePoint(source)
	hashNegativeBlock := getSourceKey(source)
	hashPositiveBlock := getSourceKey(source)
	origin := getSourceKey(source)
	prev := getSourceKey(source)
	voucherBlock := impl.NewVoucherBlock(point, hashNegativeBlock, hashPositiveBlock, origin)
	voucherBlock.SetSign(sign)
	voucherBlock.SetPreviousHash(blockpkg.FirstPrevHash, prev)
	return voucherBlock
}

func newConsumptionBlock(value []byte) interfaces.IBlock {
	source := common.NewZeroCopySource(value)
	sign := getSourceSign(source)
	var positiveBlocks []interfaces.IBlock
	m, eof := source.NextUint64()
	tryEof(eof)
	for i := 0; i < int(m); i++ {
		positiveBlocks = append(positiveBlocks, getSourcePositiveBlock(source))
	}
	prev := getSourceKey(source)
	consumptionBlock := impl.NewConsumptionBlock(positiveBlocks)
	consumptionBlock.SetSign(sign)
	consumptionBlock.SetPreviousHash(blockpkg.FirstPrevHash, prev)
	return consumptionBlock
}

func getSourcePositiveBlock(source *common.ZeroCopySource) interfaces.IBlock {
	sign := getSourceSign(source)
	point := getSourcePoint(source)
	to := getSourceKey(source)
	prev := getSourceKey(source)
	positiveBlock := impl.NewPositiveBlock(point, to)
	positiveBlock.SetSign(sign)
	positiveBlock.SetPreviousHash(blockpkg.FirstPrevHash, prev)
	return positiveBlock
}

func getSourceKey(source *common.ZeroCopySource) interfaces.IKey {
	t, eof := source.NextByte()
	tryEof(eof)
	hash, eof := source.NextHash()
	tryEof(eof)
	return header.NewKey(blockpkg.BlockType(t), hash)
}

func getSourceSign(source *common.ZeroCopySource) []byte {
	buff, _, _, eof := source.NextVarBytes()
	tryEof(eof)
	return buff
}

func getSourcePoint(source *common.ZeroCopySource) interfaces.IPoint {
	origin := getSourceKey(source)
	timeStamp, eof := source.NextUint64()
	tryEof(eof)
	nonce, eof := source.NextUint32()
	tryEof(eof)
	expireDate, eof := source.NextUint64()
	tryEof(eof)
	sign, _, _, eof := source.NextVarBytes()
	tryEof(eof)
	point := body.NewPoint(origin, timeStamp, nonce, expireDate)
	point.SetSign(sign)
	return point
}

func tryEof(eof bool) {
	if eof {
		cross.TryPanic(io.ErrUnexpectedEOF)
	}
}
