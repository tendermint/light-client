package keystore

import (
	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
)

var DefaultEncoder Encoder = encoder{}

// Encoder is used to encrypt any key with a passphrase for storage
type Encoder interface {
	Encrypt(key crypto.PrivKey, pass string) ([]byte, error)
	Decrypt(data []byte, pass string) (crypto.PrivKey, error)
}

func secret(passphrase string) []byte {
	// TODO: Sha256(Bcrypt(passphrase))
	return crypto.Sha256([]byte(passphrase))
}

type encoder struct{}

func (e encoder) Encrypt(key crypto.PrivKey, pass string) ([]byte, error) {
	s := secret(pass)
	cipher := crypto.EncryptSymmetric(key.Bytes(), s)
	return cipher, nil
}

func (e encoder) Decrypt(data []byte, pass string) (crypto.PrivKey, error) {
	s := secret(pass)
	private, err := crypto.DecryptSymmetric(data, s)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Passphrase")
	}
	key, err := crypto.PrivKeyFromBytes(private)
	return key, errors.Wrap(err, "Invalid Passphrase")
}
