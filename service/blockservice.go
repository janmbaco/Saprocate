package service

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	common2 "github.com/janmbaco/Saprocate/common"
	"github.com/janmbaco/Saprocate/core/store"
	"github.com/janmbaco/Saprocate/core/types"
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/body"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/impl"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/ontio/ontology/common"
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
func (this *BlockService) ExistsOrigin(originSign common.Uint256) bool {
	return this.store.Has(
		&header.Key{
			Type: blockpkg.Origin,
			Hash: originSign,
		})
}

func (this *BlockService) RegisterOrigin(origin *impl.Block) {
	this.verifyOrigin(origin)
	this.store.Save(origin)
}

func (this *BlockService) ReservePrevHash(block *impl.ChainLinkBlock) uint32 {
	this.verifyExistsOrigin(block)
	this.tryEnChain <- true
	block.PrevHashKey = this.store.GetLastKey()
	this.nonce = rand.Uint32()
	this.reserving <- true
	return this.nonce
}

func (this *BlockService) EnchainBlock(block *impl.ChainLinkBlock, nonce uint32) {
	var shouldTryEnChainClose bool
	this.trySaving <- true
	common2.TryFinally(
		func() {
			this.verifyNonce(nonce)
			this.saving <- true
			shouldTryEnChainClose = true
			this.verifyPevHash(block)
			this.verifySign(block)
			if block.GetType() != blockpkg.Negative {
				this.verifyPoints(block.Body.(blockpkg.PointsBody))
			}
			this.store.Save(block)
		}, func() {
			if shouldTryEnChainClose {
				<-this.tryEnChain
			}
			<-this.trySaving
		})
}

func (this *BlockService) GetSummary(origin header.Key) *types.Summary {

	result := &types.Summary{
		Owner:       &origin,
		PointsCards: make(map[header.Key]uint),
	}
	pxBlocks := this.store.GetAll(blockpkg.Positive)
	nxBlocks := this.store.GetAll(blockpkg.Negative)
	isInNxBlocks := func(pxBlock *impl.ChainLinkBlock) bool {
		result := false
		for _, nxBlock := range nxBlocks {
			nx := nxBlock.(*impl.ChainLinkBlock)
			if nx.Body.(*body.Negative).PositiveBlockKey == pxBlock.Header.Key {
				result = true
				break
			}
		}
		return result
	}

	for _, pxBlock := range pxBlocks {
		if *pxBlock.GetOrigin() == origin {
			px := pxBlock.(*impl.ChainLinkBlock)
			if !isInNxBlocks(px) {
				result.PointsCards[*px.Body.(*body.Positive).Point.Origin]++
			}
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

func (this *BlockService) verifyOrigin(origin *impl.Block) {
	pk := origin.Body.(*body.Origin).PublicKey
	this.rsaVerify(pk, origin.GetDataSigned(), origin.GetSign())
}

func (this *BlockService) verifyExistsOrigin(block *impl.ChainLinkBlock) {
	if !this.ExistsOrigin(block.GetOrigin().Hash) {
		panic(errors.New("the origin of the block not exits in the block chain"))
	}
}

func (this *BlockService) verifyPevHash(block *impl.ChainLinkBlock) {
	if *this.store.GetLastKey() != *block.PrevHashKey {
		panic(errors.New("the previous hash reserved is distinct to the previous hash of the block"))
	}
}

func (this *BlockService) verifyNonce(nonce uint32) {
	if this.nonce != nonce {
		panic(errors.New("the nonce is incorrect"))
	}
}

func (this *BlockService) verifySign(block *impl.ChainLinkBlock) {
	pk := this.getOriginBlock(block).Body.(*body.Origin).PublicKey
	this.rsaVerify(pk, block.GetDataSigned(), block.GetSign())
}

func (this *BlockService) verifyPoints(block blockpkg.PointsBody) {
	var origin *header.Key

	txBlocks := this.store.GetAll(blockpkg.Transfer)
	txFromOrigin := make([]*body.Transfer, 0)

	fillTransferFromOrigin := func(from header.Key) {
		for _, txBlock := range txBlocks {
			if *txBlock.GetOrigin() == from {
				txFromOrigin = append(txFromOrigin, txBlock.(*body.Transfer))
			}
		}
	}
	verifyPointNotTransferred := func(point *body.Point) {
		for _, tx := range txFromOrigin {
			for _, p := range tx.GetPoints() {
				if p == point {
					panic(errors.New("the point is already transferred to another person"))
				}
			}
		}
	}

	for _, point := range block.GetPoints() {
		if origin == nil {
			fillTransferFromOrigin(*point.Origin)
		} else if point.Origin != origin {
			panic(fmt.Errorf("the point must have the same origin. \n origin expected %v \n origin found: %v", origin.Hash, point.Origin.Hash))
		}
		verifyPointNotTransferred(point)
		originBlock := this.store.Get(point.Origin).(*impl.Block)
		pk := originBlock.Body.(*body.Origin).PublicKey
		this.rsaVerify(pk, point.GetDataSigned(), point.Sign)
		origin = point.Origin
	}
}

func (this *BlockService) rsaVerify(pk *rsa.PublicKey, data []byte, signature []byte) {
	digest := sha256.Sum256(data)
	cross.TryPanic(rsa.VerifyPKCS1v15(pk, crypto.SHA256, digest[:], signature))
}

func (this *BlockService) getOriginBlock(block *impl.ChainLinkBlock) *impl.Block {
	var result *impl.Block
	originBlock := this.store.Get(block.GetOrigin())
	if originBlock.GetType() == blockpkg.Positive {
		result = this.store.Get(block.GetOrigin()).(*impl.Block)
	} else {
		result = originBlock.(*impl.Block)
	}
	return result
}
