package store

import "github.com/janmbaco/Saprocate/core/types/block"

type Interface interface{
	Save(block block.Interface)
	SaveBatch(block []block.Interface)
	Get(key *block.Key) block.Interface
	GetAll(t block.Type) []block.Interface
}
