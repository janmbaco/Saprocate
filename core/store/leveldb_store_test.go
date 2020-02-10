package store

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	common2 "github.com/janmbaco/Saprocate/common"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/body"
	"math/rand"
	"os"
	"sync"
	"testing"

	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/ontio/ontology/common"
)

var testLevelDB *LevelDbStore


func TestMain(m *testing.M) {
	dbFile := "./test"
	testLevelDB = NewLevelDBStore(dbFile, common2.NewCrypter([]byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}))
	testLevelDB.Open()
	defer func(){
		testLevelDB.Close()
		os.RemoveAll(dbFile)
	} ()
	m.Run()
}

func TestSaveDB(t *testing.T) {
	var wg sync.WaitGroup

	for i := 0; i < 1000; i++ {
		h := sha256.New()
		bs := make([]byte, 4)
		binary.LittleEndian.PutUint32(bs, uint32(i))
		h.Write(bs)
		sum := h.Sum(nil)
		ui256, _ := common.Uint256ParseFromBytes(sum)
		positive := &blockpkg.ChainLinkBlock{
			Block: blockpkg.Block{
				Header: &blockpkg.Header{
					Key: &blockpkg.Key{
						Type: blockpkg.Positive,
						Hash: ui256,
					},
					Sign: bs,
				},
				Body: &body.Positive{
					Point: &blockpkg.Point{
						Origin: &blockpkg.Key{
							Type: blockpkg.Origin,
							Hash: ui256,
						},
						To: &blockpkg.Key{
							Type: blockpkg.Origin,
							Hash: common.UINT256_EMPTY,
						},
						Timestamp: 0,
						Sign:      bs,
					},
				},
			},
			PrevHashKey: &blockpkg.Key{
				Type: blockpkg.Origin,
				Hash: common.UINT256_EMPTY,
			},
		}
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
	max := 1000- 1
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
			defer func(){
				if re := recover(); re != nil {
					t.Log(fmt.Printf("Not found %v", i))
				}
			}()
			block := testLevelDB.Get(&blockpkg.Key{
				Type: blockpkg.Positive,
				Hash: ui256,
			})
			if block.(*blockpkg.ChainLinkBlock).Header.Key.Hash != ui256 {
				t.Log("incorrecto")
			}
		}(&wg)
	}
	wg.Wait()
}


func TestGetAllDB(t *testing.T) {
	arblock := testLevelDB.GetAll(blockpkg.Positive)
	t.Log(len(arblock))
}


