package store

import (
	"crypto/sha256"
	"encoding/binary"
	"github.com/janmbaco/Saprocate/core"
	"github.com/janmbaco/Saprocate/core/types"
	"github.com/janmbaco/Saprocate/core/types/block"
	"github.com/ontio/ontology/common"
	"math/rand"
	"os"
	"testing"
)

var testLevelDB *LevelDbStore
var testKey = []byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}

func TestMain(m *testing.M) {
	dbFile := "./test"

	testLevelDB = NewLevelDBStore(dbFile, core.NewCrypter(testKey))
	m.Run()
	testLevelDB.Close()
	os.RemoveAll(dbFile)
}

func TestSaveDB(t *testing.T) {
	for i := 0; i<10000; i++ {
		h := sha256.New()
		bs := make([]byte,4)
		binary.LittleEndian.PutUint32(bs, uint32(i))
		h.Write(bs)
		sum := h.Sum(nil)
		ui256, _ := common.Uint256ParseFromBytes(sum)
		positive := &types.PositiveBlock{
			Key:      block.Key{
				Type: block.Positive,
				Sign: ui256,
			},
			Previous: &block.Key{
				Type: block.Origin,
				Sign: ui256,
			},
			Coin:     &block.Coin{
				Origin:    &block.Key{
					Type: block.Origin,
					Sign: ui256,
				},
				Timestamp: 0,
				Sign:      ui256,
			},
		}
		testLevelDB.Save(positive)
	}
}



func TestGetDB(t *testing.T){
	for i := 0; i<10000; i++ {
		min := 0
		max := 10000- 1
		random := rand.Intn(max - min)+min
		h := sha256.New()
		bs := make([]byte,4)
		binary.LittleEndian.PutUint32(bs, uint32(random))
		h.Write(bs)
		sum := h.Sum(nil)
		ui256, _ := common.Uint256ParseFromBytes(sum)
		block :=  testLevelDB.Get(&block.Key{
			Type: block.Positive,
			Sign: ui256,
		})
		if block.(*types.PositiveBlock).Sign != ui256{
			t.Log("incorrecto")
		}
	}
}

func TestSaveBatchDB(t *testing.T) {
	var arBlock []block.Interface
	for i := 10000; i<20000; i++ {
		h := sha256.New()
		bs := make([]byte,4)
		binary.LittleEndian.PutUint32(bs, uint32(i))
		h.Write(bs)
		sum := h.Sum(nil)
		ui256, _ := common.Uint256ParseFromBytes(sum)
		arBlock = append(arBlock, &types.PositiveBlock{
			Key:      block.Key{
				Type: block.Positive,
				Sign: ui256,
			},
			Previous: &block.Key{
				Type: block.Origin,
				Sign: ui256,
			},
			Coin:     &block.Coin{
				Origin:    &block.Key{
					Type: block.Origin,
					Sign: ui256,
				},
				Timestamp: 0,
				Sign:      ui256,
			},
		})
	}
	testLevelDB.SaveBatch(arBlock)
}

func TestGetAllDB(t *testing.T){
	arblock := testLevelDB.GetAll(block.Positive)
	t.Log(len(arblock))
}

