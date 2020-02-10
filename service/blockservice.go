package service

import (
	"crypto"
	"crypto/rsa"
	"crypto/sha256"
	"errors"
	"fmt"
	store2 "github.com/janmbaco/Saprocate/core/store"
	"github.com/janmbaco/Saprocate/core/types"
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/body"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/ontio/ontology/common"
	"sync"
)

type (
	prevHashReserve  struct{
		origin *blockpkg.Key
		prevHash *blockpkg.Key
	}
	BlockService struct{
		m 				sync.Mutex
		prevHashReserve *prevHashReserve
		store 			store2.Interface
	})

func NewBlockService(store store2.Interface) *BlockService {
	return &BlockService{
		prevHashReserve: nil,
		store:           store,
	}
}

func(this *BlockService) ExistsOrigin(originSign common.Uint256) bool {
	return this.store.Has(
		&blockpkg.Key{
			Type: blockpkg.Origin,
			Hash: originSign,
		})
}

func(this *BlockService) RegisterOrigin(origin *blockpkg.Block)  {
	this.onLock(func() {
		this.verifyOrigin(origin)
		this.store.Save(origin)
	})
}


func(this *BlockService) ReservePrevHash(block *blockpkg.ChainLinkBlock) {
	this.onLock(func() {
		this.verifyExistsOrigin(block)
		lastKey := this.store.GetLastKey()
		block.SetPrevHash(lastKey)
		this.prevHashReserve = &prevHashReserve{
			origin:   block.GetOrigin(),
			prevHash: lastKey,
		}
	})
}

func(this *BlockService) EnchainBlock(block *blockpkg.ChainLinkBlock)  {
	this.onLock(
		func() {
			this.verifyPevHash(block)
			this.verifySign(block)
			if block.GetType() != blockpkg.Negative {
				this.verifyPoints(block.Body.(blockpkg.PointsBody))
			}
			this.store.Save(block)
		})
}

func(this *BlockService) GetSummary(origin blockpkg.Key) *types.Summary {

	result := &types.Summary{
		Owner:       &origin,
		PointsCards: make(map[blockpkg.Key]uint),
	}
	pxBlocks := this.store.GetAll(blockpkg.Positive)
	nxBlocks := this.store.GetAll(blockpkg.Negative)
	isInNxBlocks := func(pxBlock *blockpkg.ChainLinkBlock) bool {
		result := false
		for _, nxBlock := range nxBlocks{
			nx := nxBlock.(*blockpkg.ChainLinkBlock)
			if nx.Body.(*body.Negative).PositiveBlockKey == pxBlock.Header.Key{
				result = true
				break
			}
		}
		return result
	}

	for _, pxBlock := range pxBlocks{
		if *pxBlock.GetOrigin() == origin{
			px := pxBlock.(*blockpkg.ChainLinkBlock)
			if !isInNxBlocks(px) {
				result.PointsCards[*px.Body.(*body.Positive).Point.Origin]++
			}
		}
	}

	return result

}

func(this *BlockService) verifyOrigin(origin *blockpkg.Block){
	pk := origin.Body.(*body.Origin).PublicKey
	this.rsaVerify(pk, origin.GetDataSigned(), origin.GetSign())
}

func(this *BlockService) verifyExistsOrigin(block *blockpkg.ChainLinkBlock){
	if !this.ExistsOrigin(block.GetOrigin().Hash) {
		panic(errors.New("the origin of the block not exits in the block chain"))
	}
}

func(this *BlockService) verifyPevHash(block *blockpkg.ChainLinkBlock){

	if this.prevHashReserve == nil{
		panic(errors.New("there isn't any reserve of previous hash"))
	}

	if *this.prevHashReserve.origin != *block.GetOrigin(){
		panic(errors.New("there is a reserve from other origin"))
	}

	if *this.prevHashReserve.prevHash != *block.GetPrevHash(){
		panic(errors.New("the previous hash reserved is distinct to the previous hash of the block"))
	}

}

func(this *BlockService) verifySign(block *blockpkg.ChainLinkBlock){
	pk := this.getOriginBlock(block).Body.(*body.Origin).PublicKey
   	this.rsaVerify(pk, block.GetDataSigned(), block.GetSign())
}

func(this *BlockService) verifyPoints(block blockpkg.PointsBody) {
	var origin *blockpkg.Key

	txBlocks := this.store.GetAll(blockpkg.Transfer)
	txFromOrigin := make([]*body.Transfer,0)

	fillTransferFromOrigin := func(from blockpkg.Key){
		for _, txBlock := range txBlocks {
			if *txBlock.GetOrigin() == from {
				txFromOrigin = append(txFromOrigin, txBlock.(*body.Transfer))
			}
		}
	}
	verifyPointNotTransferred := func(point *blockpkg.Point)  {
		for _, tx := range txFromOrigin {
			for _, p := range tx.GetPoints() {
				if p == point {
					panic(errors.New("the point is already transferred to another person"))
				}
			}
		}
	}

	for _, point := range block.GetPoints(){
		if origin == nil{
			fillTransferFromOrigin(*point.Origin)
		} else if point.Origin != origin {
			panic(fmt.Errorf("the point must have the same origin. \n origin expected %v \n origin found: %v", origin.Hash, point.Origin.Hash))
		}
		verifyPointNotTransferred(point)
		originBlock := this.store.Get(point.Origin).(*blockpkg.Block)
		pk := originBlock.Body.(*body.Origin).PublicKey
		this.rsaVerify(pk, point.GetDataSigned(), point.Sign)
		origin = point.Origin
	}
}

func(this *BlockService) rsaVerify(pk *rsa.PublicKey, data []byte, signature []byte) {
	digest := sha256.Sum256(data)
	cross.TryPanic(rsa.VerifyPKCS1v15(pk, crypto.SHA256, digest[:], signature))
}

func(this *BlockService) getOriginBlock(block *blockpkg.ChainLinkBlock) *blockpkg.Block {
	var result *blockpkg.Block
	originBlock := this.store.Get(block.GetOrigin())
	if  originBlock.GetType() == blockpkg.Positive {
		result = this.store.Get(block.GetOrigin()).(*blockpkg.Block)
	} else{
		result = originBlock.(*blockpkg.Block)
	}
	return result
}

func (this *BlockService) onLock(callBack func())  {
	this.m.Lock()
	defer func(){
		this.m.Unlock()
		if re := recover(); re != nil {
			panic(re.(error))
		}
	}()
	callBack()
}



