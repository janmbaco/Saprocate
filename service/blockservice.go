package service

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	common2 "github.com/janmbaco/Saprocate/common"
	"github.com/janmbaco/Saprocate/core/store"
	"github.com/janmbaco/Saprocate/core/types"
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/body"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/syndtr/goleveldb/leveldb/util"
	"math/rand"
	"time"
)

const MaxWaitForReserve = 1000

type (
	BlockService struct {
		tryEnChain chan bool
		trySaving  chan bool
		reserving  chan bool
		saving     chan bool
		nonce      uint32
		store      store.Interface
	}
)

func NewBlockService(store store.Interface) *BlockService {
	blockService := &BlockService{
		tryEnChain: make(chan bool, 1),
		trySaving:  make(chan bool, 1),
		reserving:  make(chan bool),
		saving:     make(chan bool),
		store:      store,
	}
	go blockService.reservePrevHashLoop()
	return blockService
}
func (this *BlockService) ExistsOrigin(origin interfaces.IKey) bool {
	return this.store.Has(origin)
}

func (this *BlockService) RegisterOrigin(origin interfaces.IBlock) {
	this.verifyOrigin(origin)
	this.store.Save(origin)
}

func (this *BlockService) ReservePrevHash(block interfaces.IBlock) uint32 {
	this.verifyExistsOrigin(block)
	this.tryEnChain <- true
	lastkey := this.store.GetLastKey()
	prevHashType := blockpkg.FirstPrevHash
	if block.GetHeader().GetType() == blockpkg.Negative && this.store.GetType() == store.Locurum {
		prevHashType = blockpkg.SecondPrvHash
	}
	block.SetPreviousHash(prevHashType, lastkey)
	this.nonce = rand.Uint32()
	this.reserving <- true
	return this.nonce
}

func (this *BlockService) EnchainBlock(block interfaces.IBlock, nonce uint32) {
	var shouldTryEnChainClose bool
	this.trySaving <- true
	common2.TryFinally(
		func() {
			this.verifyNonce(nonce)
			this.saving <- true
			shouldTryEnChainClose = true
			this.verifyPevHash(block)
			this.verifySign(block)
			this.verifyBlock(block)
			this.store.Save(block)
		}, func() {
			if shouldTryEnChainClose {
				<-this.tryEnChain
			}
			<-this.trySaving
		})
}

func (this *BlockService) GetSummary(owner interfaces.IKey) *types.Summary {
	result := types.NewSummary(owner)
	remainingBlocks := make(map[interfaces.IKey]bool)
	positiveBlocks := this.store.Query(util.Range{
		Start: []byte{byte(blockpkg.Positive)},
		Limit: []byte{byte(blockpkg.Consumption)},
	}, func(b interfaces.IBlock) bool {
		var selected bool
		if b.GetOrigin() == owner {
			if b.GetHeader().GetType() == blockpkg.Negative {
				remainingBlocks[b.GetBody().(*body.Negative).GetPositiveBlock().GetHeader().GetKey()] = true
			} else if b.GetHeader().GetType() == blockpkg.Consumption {
				for _, positiveBlock := range b.GetBody().(*body.Consumption).GetPositiveBlocks() {
					remainingBlocks[positiveBlock.GetHeader().GetKey()] = true
				}
			} else {
				selected = true
			}
		}
		return selected
	})

	for _, pxBlock := range positiveBlocks {
		if !remainingBlocks[pxBlock.GetHeader().GetKey()] {
			result.PointsCards[pxBlock.GetBody().(*body.Positive).GetPoint().GetOrigin()]++
		}
	}
	return result
}

func (this *BlockService) reservePrevHashLoop() {
	for {
		<-this.reserving
		var timeOut bool
		select {
		case <-this.saving:
			break
		case <-time.After(MaxWaitForReserve * time.Millisecond):
			timeOut = true
			break
		}
		this.trySaving <- true
		if timeOut {
			this.nonce = rand.Uint32()
			<-this.tryEnChain
		}
		<-this.trySaving
	}
}

func (this *BlockService) verifyOrigin(origin interfaces.IBlock) {
	pk := origin.GetBody().(*body.Origin).GetPublicKey()
	this.rsaVerify(pk, origin.GetDataSigned(), origin.GetSign())
}

func (this *BlockService) verifyExistsOrigin(block interfaces.IBlock) {
	if !this.ExistsOrigin(block.GetOrigin()) {
		panic(errors.New("the origin of the block not exits in the block chain"))
	}
}

func (this *BlockService) verifyPevHash(block interfaces.IBlock) {
	prevHashType := blockpkg.FirstPrevHash

	if block.GetHeader().GetType() == blockpkg.Negative {
		if this.store.GetType() == store.Locurum {
			prevHashType = blockpkg.SecondPrvHash
			if block.GetPreviousHash(blockpkg.FirstPrevHash) == nil {
				panic(errors.New("it's necessary a previous hash from fidelis blockchain"))
			}
		} else {
			if block.GetPreviousHash(blockpkg.FirstPrevHash) == nil {
				panic(errors.New("it's necessary a previous hash from locorum blockchain"))
			}
		}
	}

	if this.store.GetLastKey() != block.GetPreviousHash(prevHashType) {
		panic(errors.New("the previous hash reserved is distinct to the previous hash of the block"))
	}
}

func (this *BlockService) verifyNonce(nonce uint32) {
	if this.nonce != nonce {
		panic(errors.New("the nonce is incorrect"))
	}
}

func (this *BlockService) verifySign(block interfaces.IBlock) {
	pk := this.getOriginBlock(block.GetOrigin()).GetBody().(*body.Origin).GetPublicKey()
	this.rsaVerify(pk, block.GetDataSigned(), block.GetSign())
}

func (this *BlockService) verifyBlock(block interfaces.IBlock) {
	switch block.GetHeader().GetType() {
	case blockpkg.Positive:
		this.verifyPositive(block)
	case blockpkg.Negative:
		this.verifyNegative(block)
	case blockpkg.Consumption:
		this.verifyConsumption(block)
	case blockpkg.Voucher:
		this.verifyVoucher(block)
	}
}

func (this *BlockService) verifyPositive(block interfaces.IBlock) {
	this.verifyPoint(block.GetBody().(*body.Positive).GetPoint())
}

func (this *BlockService) verifyNegative(block interfaces.IBlock) {
	positiveBlock := block.GetBody().(*body.Negative).GetPositiveBlock()
	if this.store.GetType() == store.Locurum && !this.store.Any(blockpkg.Voucher,
		func(block interfaces.IBlock) bool {
			voucherBody := block.GetBody().(*body.Voucher)
			return voucherBody.GetPositivekey() == positiveBlock.GetHeader().GetKey() && voucherBody.GetPoint().GetHash() == positiveBlock.GetBody().(*body.Positive).GetPoint().GetHash()
		}) {
		panic(errors.New("there isn`t any voucher of this point"))
	} else if this.store.GetType() == store.Fidelis && !this.store.Has(positiveBlock.GetHeader().GetKey()) {
		panic(errors.New("there isn`t any voucher of this point"))
	}
}

func (this *BlockService) verifyVoucher(block interfaces.IBlock) {
	voucherBody := block.GetBody().(*body.Voucher)
	if voucherBody.HasNegative() {
		if !this.store.Has(voucherBody.GetNegativeKey()) {
			panic(errors.New("there isn`t any transfer of this point"))
		} else {
			negativeOrigin := this.store.Get(voucherBody.GetNegativeKey()).GetOrigin()
			if negativeOrigin != block.GetOrigin() {
				panic(errors.New("only the owner of the point is capable to transfer it"))
			}
			if this.store.GetType() == store.Fidelis && !this.store.Has(voucherBody.GetPositivekey()) {
				panic(errors.New("the voucher is only stored for locorum o fidelis that has the positive block"))
			}
		}
	} else if block.GetOrigin() != voucherBody.GetPoint().GetOrigin() {
		panic(errors.New("only the origin of this point is capable to give this point"))
	}
}

func (this *BlockService) verifyConsumption(block interfaces.IBlock) {
	payment := block.GetBody().(*body.Consumption)
	var pointOrigin interfaces.IKey
	for _, positiveBlock := range payment.GetPositiveBlocks() {
		if positiveBlock.GetHeader().GetType() != blockpkg.Positive {
			panic(errors.New("the payment is only possible with positive blocks"))
		}
		this.verifySign(positiveBlock)
		this.verifyBlock(positiveBlock)
		positiveBody := positiveBlock.GetBody().(*body.Positive)
		if this.store.Any(blockpkg.Negative, func(block interfaces.IBlock) bool {
			return positiveBlock.GetHeader().GetHash() == block.GetBody().(*body.Negative).GetPositiveBlock().GetHeader().GetHash()
		}) {
			panic(errors.New("the payment is impossible with points transferred"))
		} else if !this.store.Any(blockpkg.Voucher, func(block interfaces.IBlock) bool {
			voucherBody := block.GetBody().(*body.Voucher)
			return positiveBlock.GetHeader().GetKey() == voucherBody.GetPositivekey() && positiveBody.GetPoint().GetHash() == voucherBody.GetPoint().GetHash()
		}) {
			panic(errors.New("there isn't any voucher of this positive block"))
		} else if pointOrigin != nil && positiveBody.GetPoint().GetOrigin() != pointOrigin {
			panic(errors.New("there are several origin of the points, in a payment all point's origin must to be the same"))
		}
	}
}

func (this *BlockService) verifyPoint(point interfaces.IPoint) {
	this.rsaVerify(this.getOriginBlock(point.GetOrigin()).GetBody().(*body.Origin).GetPublicKey(), point.GetDataSigned(), point.GetSign())
}

func (this *BlockService) getOriginBlock(origin interfaces.IKey) interfaces.IBlock {
	originBlock := this.store.Get(origin)
	if originBlock.GetHeader().GetType() != blockpkg.Origin {
		panic(errors.New("the origin block set is not a origin block type"))
	}
	return originBlock
}

func (this *BlockService) rsaVerify(pk *rsa.PublicKey, data []byte, signature []byte) {
	digest := sha256.Sum256(data)
	cross.TryPanic(rsa.VerifyPKCS1v15(pk, crypto.SHA256, digest[:], signature))
}
