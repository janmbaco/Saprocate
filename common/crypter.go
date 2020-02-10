package common

import (
	"crypto/aes"
	"crypto/cipher"
	"github.com/janmbaco/go-reverseproxy-ssl/cross"
)

type Crypter struct{
	aead cipher.AEAD
}

func NewCrypter(key []byte) *Crypter{
	block, err := aes.NewCipher(key)
	cross.TryPanic(err)
	aead, err := cipher.NewGCM(block)
	cross.TryPanic(err)
	return &Crypter{
		aead : aead,
	}
}

func(this *Crypter) Encrypt(value []byte) []byte {
	nonce := make([]byte, this.aead.NonceSize())
	return this.aead.Seal(nonce, nonce, value, nil)
}

func (this *Crypter) Decrypt(value []byte) []byte {
	nonceSize := this.aead.NonceSize()
	nonce, cipherValue := value[:nonceSize], value[nonceSize:]
	plainValue, err := this.aead.Open(nil, nonce, cipherValue, nil)
	cross.TryPanic(err)
	return plainValue
}
