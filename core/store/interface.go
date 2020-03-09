package store

import (
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/syndtr/goleveldb/leveldb/util"
)

type Interface interface {
	Has(interfaces.IKey) bool
	Save(interfaces.IBlock)
	Get(interfaces.IKey) interfaces.IBlock
	GetType() StoreType
	Any(blockpkg.BlockType, func(block interfaces.IBlock) bool) bool
	Query(util.Range, func(interfaces.IBlock) bool) []interfaces.IBlock
	GetWhere(blockpkg.BlockType, func(block interfaces.IBlock) bool) []interfaces.IBlock
	GetAll(blockpkg.BlockType) []interfaces.IBlock
	GetLastKey() interfaces.IKey
	Open()
	Close()
}
