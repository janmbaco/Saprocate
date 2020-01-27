package store

import (
	"github.com/ethereum/go-ethereum/common/fdlimit"
	"github.com/janmbaco/Saprocate/core"
	"github.com/janmbaco/Saprocate/core/types"
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
	"github.com/syndtr/goleveldb/leveldb"
	"github.com/syndtr/goleveldb/leveldb/filter"
	"github.com/syndtr/goleveldb/leveldb/opt"
	"github.com/syndtr/goleveldb/leveldb/util"
)

const BITSPERKEY = 10


type LevelDbStore struct{
	db *leveldb.DB
	crypter *core.Crypter
	batch *leveldb.Batch
}

func NewLevelDBStore(file string, crypter *core.Crypter) *LevelDbStore{
	openFileCache := opt.DefaultOpenFilesCacheCapacity
	maxOpenFiles, err := fdlimit.Current()
	cross.TryPanic(err)
	if maxOpenFiles < openFileCache*5{
		openFileCache = maxOpenFiles / 5
	}
	
	if openFileCache < 16 {
		openFileCache = 16
	}
	
	options := opt.Options{
		Filter:	filter.NewBloomFilter(BITSPERKEY),
		OpenFilesCacheCapacity:	openFileCache,
	}

	db, err := leveldb.OpenFile(file, &options)
	cross.TryPanic(err)

	return &LevelDbStore{
		db:    db,
		crypter: crypter,
		batch: nil,
	}
}

func(this *LevelDbStore) Save(block block.Interface) error {
	return this.db.Put(block.SerializeKey(), this.crypter.Encrypt(block.SerializeValue()) , nil)
}

func(this *LevelDbStore) SaveBatch(blocks []block.Interface) error {
	this.batch = new(leveldb.Batch)
	for _, block := range blocks{
		this.batch.Put(block.SerializeKey(), this.crypter.Encrypt(block.SerializeValue()))
	}
	return this.db.Write(this.batch, nil)
}

func(this *LevelDbStore) Get(key *block.Key) block.Interface {
	dat, err := this.db.Get(key.SerializeKey(), nil)
	cross.TryPanic(err)
	return types.NewFromBytes(key, this.crypter.Decrypt(dat))

}

func(this *LevelDbStore) GetAll(t block.Type) []block.Interface {
	var result []block.Interface
	iter := this.db.NewIterator(util.BytesPrefix([]byte{byte(t)}), nil)
	for iter.Next(){
		result = append(result, types.NewFromBytes(types.KeyFromBytes(iter.Key()), this.crypter.Decrypt(iter.Value())))
	}
	iter.Release()
	err := iter.Error()
	cross.TryPanic(err)
	return result
}

func (this *LevelDbStore) Close() {
	err := this.db.Close()
	cross.TryPanic(err)
}



