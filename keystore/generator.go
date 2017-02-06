package keystore

import crypto "github.com/tendermint/go-crypto"

// Generator determines the type of Private Key we make by default
type Generator interface {
	Generate() crypto.PrivKey
}

type GenFunc func() crypto.PrivKey

func (f GenFunc) Generate() crypto.PrivKey {
	return f()
}

var (
	GenEd25519 GenFunc = func() crypto.PrivKey {
		return crypto.GenPrivKeyEd25519()
	}
	GenSecp256k1 GenFunc = func() crypto.PrivKey {
		return crypto.GenPrivKeySecp256k1()
	}
)
