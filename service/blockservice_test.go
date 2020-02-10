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
	"github.com/ontio/ontology/common"
	"os"
	"sync"
	"testing"
	"time"
)

var blockService *BlockService

type keyPairStruct struct{
	privateKey 	*rsa.PrivateKey
	publicKey *rsa.PublicKey
	key *blockpkg.Key
	sign []byte
}
var keyPair []*keyPairStruct



func TestMain(m *testing.M) {
	dbFile := "./test"
	testLevelDB := store2.NewLevelDBStore(dbFile, common2.NewCrypter([]byte{0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7, 0x8, 0x9, 0xa, 0xb, 0xc, 0xd, 0xe, 0xf}))
	testLevelDB.Open()
	keyPair =  make([]*keyPairStruct,4)
	for i := 0; i< 4; i++{
		keyPair[i], _ = generateKeyPair(2048)
		block := &blockpkg.Block{
			Header: nil,
			Body:   &body.Origin{PublicKey: keyPair[i].publicKey},
		}
		sign, _ := sign(keyPair[i].privateKey, block.GetDataSigned())
		signSum := sha256.Sum256(sign)
		ui256, _ := common.Uint256ParseFromBytes(signSum[:])
		keyPair[i].key = &blockpkg.Key{
			Type: blockpkg.Origin,
			Hash: ui256,
		}
		keyPair[i].sign = sign
	}
	blockService = NewBlockService(testLevelDB)
	defer func(){
		testLevelDB.Close()
		os.RemoveAll(dbFile)
	} ()
	m.Run()
}

func TestRegisterOrigins(t *testing.T){
	var wg sync.WaitGroup
	for _, keypair := range keyPair{
		wg.Add(1)
		go func(keypair *keyPairStruct){
			defer wg.Done()
			blockService.RegisterOrigin(&blockpkg.Block{
				Header: &blockpkg.Header{
					Key:  keypair.key,
					Sign: keypair.sign,
				},
				Body:   &body.Origin{
					PublicKey:keypair.publicKey,
				},
			})

		}(keypair)

	}
	wg.Wait()
}

func TestGivePoints(t *testing.T){
	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		for j := 0; j<2;j++{
			wg.Add(1)
			go func(j int){
				defer wg.Done()
				point := blockpkg.Point{
					Origin:    keyPair[j].key,
					To:        keyPair[j+2].key,
					Timestamp: uint64(time.Now().UnixNano()),
					Sign:      nil,
				}
				point.Sign, _ = sign(keyPair[j].privateKey, point.GetDataSigned())
				positive := &blockpkg.ChainLinkBlock{
					Block:       blockpkg.Block{
						Header: &blockpkg.Header{
							Key:  keyPair[j].key,
							Sign: nil,
						},
						Body:   &body.Positive{
							Point:&point,
						},
					},
					PrevHashKey: nil,
				}
				blockService.ReservePrevHash(positive)
				positive.Header.Sign, _ = sign(keyPair[j+2].privateKey, positive.GetDataSigned())
				blockService.EnchainBlock(positive)

			}(j)

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
