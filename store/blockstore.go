package store

import (
	"github.com/janmbaco/Saprocate/core/types"
	"github.com/ontio/ontology/common"
)

type BlockStore struct {

}

func(this *BlockStore) Put(key *common.Uint256, data []byte){
	//TODO Put data store
}

func(this *BlockStore) GetCurrency(key *common.Uint256) (*types.Currency, error) {
	//TODO get currency
	return nil, nil
}

func(this *BlockStore) GetTransaction(key *common.Uint256)(*types.Transaction, error){
	//TODO get transaction
	return nil, nil
}

