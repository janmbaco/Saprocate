package store

import "github.com/janmbaco/Saprocate/core/types/blockpkg"

type Interface interface {
	Has(key *blockpkg.Key) bool
	Save(block blockpkg.Interface)
	Get(key *blockpkg.Key) blockpkg.Interface
	GetAll(t blockpkg.Type) []blockpkg.Interface
	GetLastKey() *blockpkg.Key
	Open()
	Close()
}
