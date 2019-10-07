package services

import (
	"crypto/rsa"
	"encoding/asn1"
	"sync"

	"github.com/golang/protobuf/proto"
	"github.com/janmbaco/Saprocate/core/memory"
	"github.com/janmbaco/Saprocate/store"
	pb "github.com/janmbaco/saprocate/core/types/protobuf"
	"github.com/ontio/ontology/common"
)

type prevHashReserve struct {
	SignActor *common.Uint256
	PrevHash  *common.Uint256
}

type blockService struct {
	prevHashReserved *prevHashReserve
	m                *sync.Mutex
	store            *store.BlockStore
	state            *memory.State
}

var BlockService *blockService

func Init() {
	BlockService = &blockService{}
}

func (this *blockService) ReservePrevHash(prevHashRequest *pb.ReservePrevHashRequest) *common.Uint256 {

	var prevHash *common.Uint256
	this.onLock(func(){

		if this.prevHashReserved != nil {
			//TODO: launch error
		}

		this.prevHashReserved.PrevHash = this.state.GetLastHash()

		prevHash =  this.prevHashReserved.PrevHash
	})
	return prevHash
}

func (this *blockService) CreateCurrency(currencyRequest *pb.CreateCurrencyRequest) {

	this.onLock(func(){

		this.checkReservedPrevHash(currencyRequest.Data.Header)

		publicKey := &rsa.PublicKey{}
		_, err := asn1.Unmarshal(currencyRequest.Data.Body.Asn1, publicKey)
		this.tryError(err)

		dataBytes, err := proto.Marshal(currencyRequest.Data)
		this.tryError(err)

		this.verifyDataSign(publicKey, currencyRequest.Hash, dataBytes)

		this.store.Put(this.prevHashReserved.SignActor, dataBytes)
	})
}

func (this *blockService) CreateTransaction(transactionRequest *pb.CreateTransactionRequest) {
	this.onLock(func(){

		this.checkReservedPrevHash(transactionRequest.Data.Header)

		dataBytes, err := proto.Marshal(transactionRequest.Data)
		this.tryError(err)

		this.verifyDataSign(this.getPublicKey(this.prevHashReserved.SignActor), transactionRequest.Hash, dataBytes)

		for _, move := range transactionRequest.Data.Body.BlockHash{
			blockHash, err := common.Uint256ParseFromBytes(move)
			this.tryError(err)
			lastReceiver := this.state.GetOwner(&blockHash)
			if lastReceiver != this.prevHashReserved.SignActor{
				//Todo launch error
			}
		}
	})
}



func(this *blockService) checkReservedPrevHash(header *pb.Header) {
	if this.prevHashReserved == nil {
		//TODO: launch error
	}

	reservedPrevHash := this.parseFromBytes(header)
	if !(*reservedPrevHash  == *this.prevHashReserved) {
		//TODO: launch error
	}

	this.verifyPrevHash(reservedPrevHash.PrevHash)
}

func (this *blockService) onLock(callBack func()){
	this.m.Lock()
	defer this.showError()
	callBack()
}

func (this *blockService) showError() {
	if re := recover(); re != nil {

	}
	this.m.Unlock()
}

func (this *blockService) tryError(err error){
	if err != nil {
		//TODO manage defined error
		panic(err)
	}
}

func (this *blockService) verifyPrevHash(prevHash *common.Uint256) {

}

func (this *blockService) getPublicKey(signActor *common.Uint256) *rsa.PublicKey {
	//TODO get public key
	return &rsa.PublicKey{}
}

func (this *blockService) verifyDataSign(publicKey *rsa.PublicKey, hash []byte, data []byte) {

}

func (this *blockService) parseFromBytes(header *pb.Header) *prevHashReserve {
	signActor, err := common.Uint256ParseFromBytes(header.SignActor)
	this.tryError(err)
	prevHash, err := common.Uint256ParseFromBytes(header.PrevHash)
	this.tryError(err)

	return &prevHashReserve{
		SignActor: &signActor,
		PrevHash:  &prevHash,
	}
}

