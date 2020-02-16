package service

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	common2 "github.com/janmbaco/Saprocate/common"
	memory2 "github.com/janmbaco/Saprocate/core/memory"
	store2 "github.com/janmbaco/Saprocate/core/store"
	"github.com/janmbaco/Saprocate/core/types"
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/Point"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/body"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/impl"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/ontio/ontology/common"
	errors2 "github.com/ontio/ontology/errors"
	"github.com/syndtr/goleveldb/leveldb/util"
	"math/rand"
	"time"
)

const MaxWaitForReserve = 10000

type (
	BlockService struct {
		tryEnChain chan bool
		trySaving  chan bool
		reserving  chan bool
		saving     chan bool
		nonce      uint32
		memory     *memory2.Memory
		store      store2.Interface
	}
)

func NewBlockService(store store2.Interface) *BlockService {
	blockService := &BlockService{
		tryEnChain: make(chan bool, 1),
		trySaving:  make(chan bool, 1),
		reserving:  make(chan bool),
		saving:     make(chan bool),
		memory:     memory2.NewMemory(),
		store:      store,
	}
	blockService.loadMemory()
	go blockService.reservePrevHashLoop()
	return blockService
}
func (this *BlockService) ExistsOrigin(originSign common.Uint256) bool {
	return this.store.Has(
		&header.Key{
			Type: header.Origin,
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
	this.trySaving <- true
	var shouldEndTryEnchain bool
	common2.TryFinally(
		func() {
			this.verifyNonce(nonce)
			this.saving <- true
			shouldEndTryEnchain = true
			this.verifyPevHash(block)
			this.verifySign(block)
			this.verifyPoints(block)
			this.store.Save(block)
			this.addToMemory(block)
		}, func() {
			if shouldEndTryEnchain {
				<- this.tryEnChain
			}
			<-this.trySaving
		})
}

func (this *BlockService) GetSummary(origin *header.Key) *types.Summary {

	result := &types.Summary{
		Owner:       origin,
		PointsCards: make(map[header.Key]uint),
	}
	pxBlocks := this.store.Query(util.BytesPrefix([]byte{byte(header.Numu)}),
		func(b blockpkg.Interface) bool {
			var result bool
			if *b.GetOrigin() == *origin {
				result = true
			}
			return result
		})

	for _, pxBlock := range pxBlocks {
		px := pxBlock.(*impl.ChainLinkBlock)
		for _, point := range px.Body.(body.PointsBody).GetPoints() {
			signSum := sha256.Sum256(point.GetSign())
			ui256, _ := common.Uint256ParseFromBytes(signSum[:])
			if this.memory.Point[ui256] == memory2.Stored {
				result.PointsCards[*px.Body.(*body.Numu).Point.GetOrigin()]++
			}
		}
	}

	return result
}

func (this *BlockService) loadMemory() {
	for _, b := range this.store.Query(&util.Range{
		Start: []byte{byte(header.Numu)},
		Limit: []byte{byte(header.Pay)},
	}, func(b blockpkg.Interface) bool {
		return true
	}) {
		bc := b.(*impl.ChainLinkBlock)
		pb := bc.Body.(body.PointsBody)
		for _, p := range pb.GetPoints() {
			this.memory.SetPoint(b.GetType(), p)
		}
	}
}

func (this *BlockService) addToMemory(block *impl.ChainLinkBlock) {
	for _, point := range block.Body.(body.PointsBody).GetPoints() {
		this.memory.SetPoint(block.GetType(), point)
	}
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

func (this *BlockService) verifyPoints(block *impl.ChainLinkBlock) {
	var origin *header.Key
	for _, point := range block.Body.(body.PointsBody).GetPoints() {
		if origin != nil && point.GetOrigin() != origin {
			panic(fmt.Errorf("the point must have the same origin. \n origin expected %v \n origin found: %v", origin.Hash, point.GetOrigin().Hash))
		}
		this.verifyStatePoint(block.GetType(), point)
		originBlock := this.store.Get(point.GetOrigin()).(*impl.Block)
		pk := originBlock.Body.(*body.Origin).PublicKey
		this.rsaVerify(pk, point.GetDataSigned(), point.GetSign())
		origin = point.GetOrigin()
	}
}

func (this *BlockService) verifyStatePoint(t header.Type, point Point.Interface) {
	signSum := sha256.Sum256(point.GetSign())
	ui256, _ := common.Uint256ParseFromBytes(signSum[:])
	state := this.memory.Point[ui256]
	switch t {
	case header.Numu:
		if state != memory2.None {
			panic(errors2.NewErr("the point is already stored"))
		}
	case header.Transfer, header.Pay:
		if state != memory2.Stored {
			panic(errors2.NewErr("this point is already spent"))
		}
	}
}

func (this *BlockService) rsaVerify(pk *rsa.PublicKey, data []byte, signature []byte) {
	digest := sha256.Sum256(data)
	cross.TryPanic(rsa.VerifyPKCS1v15(pk, crypto.SHA256, digest[:], signature))
}

func (this *BlockService) getOriginBlock(block *impl.ChainLinkBlock) *impl.Block {
	var result *impl.Block
	originBlock := this.store.Get(block.GetOrigin())
	if originBlock.GetType() == header.Numu {
		result = this.store.Get(block.GetOrigin()).(*impl.Block)
	} else {
		result = originBlock.(*impl.Block)
	}
	return result
}
