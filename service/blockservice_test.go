package service

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	common2 "github.com/janmbaco/Saprocate/common"
	store2 "github.com/janmbaco/Saprocate/core/store"
	"github.com/janmbaco/Saprocate/core/types/blockpkg"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/body"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/header"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/impl"
	"github.com/janmbaco/Saprocate/core/types/blockpkg/interfaces"
	"github.com/ontio/ontology/common"
	rand2 "math/rand"
	"os"
	"sync"
	"testing"
	"time"
)

var blockServices []*BlockService

type keyPairStruct struct {
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
	block      interfaces.IBlock
}

var keyPair []*keyPairStruct

func TestMain(m *testing.M) {
	dbFile := "./test"
	crypter := common2.NewCrypter([]byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf})
	testLevelDB := store2.NewLevelDBStore(dbFile, crypter)
	testLevelDB.Open()
	keyPair = make([]*keyPairStruct, 4)
	blockServices = make([]*BlockService, 3)
	blockServices[0] = NewBlockService(Locorum, "./locorum")
	for i := 0; i < 4; i++ {

		keyPair[i], _ = generateKeyPair(2048)
		keyPair[i].block = impl.NewOriginBlock(keyPair[i].publicKey)
		sign, _ := sign(keyPair[i].privateKey, keyPair[i].block.GetDataSigned())
		keyPair[i].block.SetSign(sign)
	}
	blockService = NewBlockService(testLevelDB)
	defer func() {
		testLevelDB.Close()
		_ = os.RemoveAll(dbFile)
	}()
	m.Run()
}

func TestRegisterOrigins(t *testing.T) {
	var wg sync.WaitGroup
	for _, keypair := range keyPair {
		wg.Add(1)
		go func(keypair *keyPairStruct) {
			defer wg.Done()
			blockService.RegisterOrigin(keypair.block)

		}(keypair)

	}
	wg.Wait()
}

func TestGivePoints(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		for j := 0; j < 2; j++ {
			wg.Add(1)
			go func(i int, j int) {
				defer wg.Done()
				point := body.NewPoint(keyPair[j].block.GetOrigin(), uint64(time.Now().UnixNano()), rand2.Uint32(), uint64(time.Now().UnixNano()))
				pointSign, _ := sign(keyPair[j].privateKey, point.GetDataSigned())
				point.SetSign(pointSign)
				positive := impl.NewPositiveBlock(point, keyPair[j+2].block.GetOrigin())

				common2.TryError(func() {
					nonce := blockServices[0].ReservePrevHash(positive)
					blockSign, _ := sign(keyPair[j+2].privateKey, positive.GetDataSigned())
					positive.SetSign(blockSign)
					blockServices[0].EnchainBlock(positive, nonce)
				}, func(err error) {
					t.Logf("i %d j %d error %s", i, j, err.Error())
				})
			}(i, j)
		}
	}
	wg.Wait()
}

// GenerateKeyPair generates a new key pair
func generateKeyPair(bits int) (*keyPairStruct, error) {
	privkey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}
	return &keyPairStruct{
		privateKey: privkey,
		publicKey:  &privkey.PublicKey,
	}, nil
}

func sign(privateKey *rsa.PrivateKey, data []byte) ([]byte, error) {
	digest := sha256.Sum256(data)
	return rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA256, digest[:])

}
