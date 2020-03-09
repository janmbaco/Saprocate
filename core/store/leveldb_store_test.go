package store

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	common2 "github.com/janmbaco/Saprocate/common"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/body"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/impl"
	"math/rand"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/ontio/ontology/common"
)

var testLevelDB *LevelDbStore

func TestMain(m *testing.M) {
	dbFile := "./test"
	testLevelDB = NewLevelDBStore(Fidelis, dbFile, common2.NewCrypter([]byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}))
	testLevelDB.Open()
	defer func() {
		testLevelDB.Close()
		os.RemoveAll(dbFile)
	}()
	m.Run()
}

func TestSaveDB(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(i))
		header := header.NewHeader(blockpkg.Origin, bs)
		point := body.NewPoint(header.GetKey(), uint64(time.Now().UnixNano()), rand.Uint32(), uint64(time.Now().UnixNano()))
		point.SetSign(bs)
		positive := impl.NewPositiveBlock(point, header.GetKey())
		positive.SetSign(bs)
		positive.SetPreviousHash(blockpkg.FirstPrevHash, testLevelDB.GetLastKey())
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			testLevelDB.Save(positive)
		}(&wg)
	}
	wg.Wait()
}

func TestGetDB(t *testing.T) {
	min := 0
	max := 1000 - 1
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		random := rand.Intn(max-min) + min
		h := sha256.New()
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(random))
		h.Write(bs)
		sum := h.Sum(nil)
		wg.Add(1)
		go func(wg *sync.WaitGroup) {
			defer wg.Done()
			ui256, _ := common.Uint256ParseFromBytes(sum)
			defer func() {
				if re := recover(); re != nil {
					t.Log(fmt.Printf("Not found %v", i))
				}
			}()
			key := header.NewKey(blockpkg.Positive, ui256)
			block := testLevelDB.Get(key)
			if block.GetHeader().GetHash() != ui256 {
				t.Log("incorrecto")
			} else {
				hash := block.GetPreviousHash(blockpkg.FirstPrevHash).GetHash()
				t.Logf(hash.ToHexString())
			}
		}(&wg)
	}
	wg.Wait()
}

func TestGetAllDB(t *testing.T) {
	arblock := testLevelDB.GetAll(blockpkg.Positive)
	t.Log(len(arblock))
}
