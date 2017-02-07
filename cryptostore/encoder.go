package cryptostore

import (
	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
)

var (
	SecretBox Encoder = secretbox{}
	Noop      Encoder = noop{}
)

// Encoder is used to encrypt any key with a passphrase for storage
type Encoder interface {
	Encrypt(key crypto.PrivKey, pass string) ([]byte, error)
	Decrypt(data []byte, pass string) (crypto.PrivKey, error)
}

func secret(passphrase string) []byte {
	// TODO: Sha256(Bcrypt(passphrase))
	return crypto.Sha256([]byte(passphrase))
}

type secretbox struct{}

func (e secretbox) Encrypt(key crypto.PrivKey, pass string) ([]byte, error) {
	s := secret(pass)
	cipher := crypto.EncryptSymmetric(key.Bytes(), s)
	return cipher, nil
}

func (e secretbox) Decrypt(data []byte, pass string) (crypto.PrivKey, error) {
	s := secret(pass)
	private, err := crypto.DecryptSymmetric(data, s)
	if err != nil {
		return nil, errors.Wrap(err, "Invalid Passphrase")
	}
	key, err := crypto.PrivKeyFromBytes(private)
	return key, errors.Wrap(err, "Invalid Passphrase")
}

type noop struct{}

func (n noop) Encrypt(key crypto.PrivKey, pass string) ([]byte, error) {
	return key.Bytes(), nil
}

func (n noop) Decrypt(data []byte, pass string) (crypto.PrivKey, error) {
	return crypto.PrivKeyFromBytes(data)
}
