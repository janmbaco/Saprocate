package types

import "crypto/rsa"

type Currency struct {
	PublicKey *rsa.PublicKey
	Purpose string
}
