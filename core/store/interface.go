package store

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Interface interface {
	Has(key *header.Key) bool
	Save(block blockpkg.Interface)
	Get(key *header.Key) blockpkg.Interface
	Query(rang *util.Range, where func(blockpkg.Interface) bool) []blockpkg.Interface
	GetAll(t blockpkg.Type) []blockpkg.Interface
	GetLastKey() *header.Key
	Open()
	Close()
}
