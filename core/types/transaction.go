package types

import "github.com/ontio/ontology/common"

type Transaction struct {
	ReceiverActor *common.Uint256
	Moves []*common.Uint256
}
