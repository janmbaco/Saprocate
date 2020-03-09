package impl

import (
	"crypto/rsa"
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/body"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
)

func NewOriginBlock(publicKey *rsa.PublicKey) interfaces.IBlock {
	originBlock := &Block{}
	originBlock.header = header.NewHeader(blockpkg.Origin, nil)
	originBlock.body = body.NewOrigin(publicKey)
	return originBlock
}

func NewPositiveBlock(point interfaces.IPoint, to interfaces.IKey) interfaces.IBlock {
	positiveBlock := &ChainLinkBlock{}
	positiveBlock.header = header.NewHeader(blockpkg.Positive, nil)
	positiveBlock.body = body.NewPositive(point, to)
	return positiveBlock
}

func NewNegativeblock(positiveBlock interfaces.IBlock) interfaces.IBlock {
	negativeBlock := &DoubleChainLinkBlock{}
	negativeBlock.header = header.NewHeader(blockpkg.Negative, nil)
	negativeBlock.body = body.NewNegative(positiveBlock)
	return negativeBlock
}

func NewVoucherBlock(point interfaces.IPoint, hashNegative interfaces.IKey, hashPositive interfaces.IKey, origin interfaces.IKey) interfaces.IBlock {
	voucherBlock := &ChainLinkBlock{}
	voucherBlock.header = header.NewHeader(blockpkg.Voucher, nil)
	voucherBlock.body = body.NewVoucher(point, hashNegative, hashPositive, origin)
	return voucherBlock
}

func NewConsumptionBlock(positivesBlock []interfaces.IBlock) interfaces.IBlock {
	paymentBlock := &ChainLinkBlock{}
	paymentBlock.header = header.NewHeader(blockpkg.Consumption, nil)
	paymentBlock.body = body.NewConsumption(positivesBlock)
	return paymentBlock
}
