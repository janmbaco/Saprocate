package store

import (
	"github.com/ethereum/go-ethereum/common/fdlimit"
	"github.com/janmbaco/Saprocate/common"
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	common2 "github.com/ontio/ontology/common"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const BITSPERKEY = 10

var cLastkey = []byte("LAST_KEY")

type LevelDbStore struct {
	storeType StoreType
	db        *leveldb.DB
	file      string
	options   *opt.Options
	crypter   *common.Crypter
	batch     *leveldb.Batch
}

func NewLevelDBStore(storeType StoreType, file string, crypter *common.Crypter) *LevelDbStore {
	openFileCache := opt.DefaultOpenFilesCacheCapacity
	maxOpenFiles, err := fdlimit.Current()
	cross.TryPanic(err)
	if maxOpenFiles < openFileCache*5 {
		openFileCache = maxOpenFiles / 5
	}

	if openFileCache < 16 {
		openFileCache = 16
	}

	options := opt.Options{
		Filter:                 filter.NewBloomFilter(BITSPERKEY),
		OpenFilesCacheCapacity: openFileCache,
	}
	return &LevelDbStore{
		storeType: storeType,
		file:      file,
		options:   &options,
		crypter:   crypter,
		batch:     nil,
	}
}

func (this *LevelDbStore) GetType() StoreType {
	return this.storeType
}

func (this *LevelDbStore) Save(b interfaces.IBlock) {
	err := this.db.Put(b.KeyToBytes(), this.crypter.Encrypt(b.ValueToBytes()), nil)
	cross.TryPanic(err)
	if b.GetHeader().GetType() != blockpkg.Origin {
		err = this.db.Put(cLastkey, this.crypter.Encrypt(b.KeyToBytes()), nil)
		cross.TryPanic(err)
	}
}

func (this *LevelDbStore) Has(key *header.Key) bool {
	b, err := this.db.Has(key.ToBytes(), nil)
	cross.TryPanic(err)
	return b
}

func (this *LevelDbStore) Get(key *header.Key) interfaces.IBlock {
	dat, err := this.db.Get(key.ToBytes(), nil)
	cross.TryPanic(err)
	return BlockFromBytes(key, this.crypter.Decrypt(dat))
}

func (this *LevelDbStore) Query(rang *util.Range, where func(block interfaces.IBlock) bool) []interfaces.IBlock {
	var result []interfaces.IBlock
	iter := this.db.NewIterator(rang, nil)
	for iter.Next() {
		block := BlockFromBytes(KeyFromBytes(iter.Key()), this.crypter.Decrypt(iter.Value()))
		if where(block) {
			result = append(result, block)
		}
	}
	iter.Release()
	err := iter.Error()
	cross.TryPanic(err)
	return result
}

func (this *LevelDbStore) GetWhere(t blockpkg.BlockType, where func(block interfaces.IBlock) bool) []interfaces.IBlock {
	return this.Query(&util.Range{Start: []byte{byte(t)}, Limit: []byte{byte(t)}}, where)
}

func (this *LevelDbStore) Any(t blockpkg.BlockType, where func(block interfaces.IBlock) bool) bool {
	result := false
	iter := this.db.NewIterator(util.BytesPrefix([]byte{byte(t)}), nil)
	for result || iter.Next() {
		block := BlockFromBytes(KeyFromBytes(iter.Key()), this.crypter.Decrypt(iter.Value()))
		result = where(block)
	}
	iter.Release()
	err := iter.Error()
	cross.TryPanic(err)
	return result
}

func (this *LevelDbStore) GetAll(t blockpkg.BlockType) []interfaces.IBlock {
	var result []interfaces.IBlock
	iter := this.db.NewIterator(util.BytesPrefix([]byte{byte(t)}), nil)
	for iter.Next() {
		result = append(result, BlockFromBytes(KeyFromBytes(iter.Key()), this.crypter.Decrypt(iter.Value())))
	}
	iter.Release()
	err := iter.Error()
	cross.TryPanic(err)
	return result
}

func (this *LevelDbStore) GetLastKey() interfaces.IKey {
	var result interfaces.IKey
	result = header.NewKey(blockpkg.Origin, common2.UINT256_EMPTY)
	b, err := this.db.Has(cLastkey, nil)
	cross.TryPanic(err)
	if b {
		dat, err := this.db.Get(cLastkey, nil)
		cross.TryPanic(err)
		result = KeyFromBytes(this.crypter.Decrypt(dat))
	}
	return result
}

func (this *LevelDbStore) Open() {
	db, err := leveldb.OpenFile(this.file, this.options)
	cross.TryPanic(err)
	this.db = db
}

func (this *LevelDbStore) Close() {
	err := this.db.Close()
	cross.TryPanic(err)
}
