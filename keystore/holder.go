package keystore

import (
	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
	lightclient "github.com/tendermint/light-client"
)

// storer combines all these things to implement lightclient.KeyStore
type storer struct {
	gen Generator
	es  encryptedStorage
}

func New(gen Generator, coder Encoder, store Storage) lightclient.KeyStore {
	return storer{
		gen: gen,
		es: encryptedStorage{
			coder: coder,
			store: store,
		},
	}
}

func (s storer) Create(name, passphrase string) error {
	key := s.gen.Generate()
	return s.es.Put(name, passphrase, key)
}

func (s storer) List() ([]lightclient.KeyInfo, error) {
	return s.es.List()
}

func (s storer) Get(name string) (lightclient.KeyInfo, error) {
	_, info, err := s.es.store.Get(name)
	return info, err
}

func (s storer) Signature(name, passphrase string, data []byte) ([]byte, error) {
	key, _, err := s.es.Get(name, passphrase)
	if err != nil {
		return nil, err
	}

	sig := key.Sign(data)
	return sig.Bytes(), nil
}

// Verify includes hardcoded byte parsing
func (s storer) Verify(data, sigBytes, pubkey []byte) error {
	sig, err := crypto.SignatureFromBytes(sigBytes)
	if err != nil {
		return errors.Wrap(err, "Verify")
	}

	pk, err := crypto.PubKeyFromBytes(pubkey)
	if err != nil {
		return errors.Wrap(err, "Verify")
	}

	valid := pk.VerifyBytes(data, sig)
	if !valid {
		return errors.New("Invalid Signature")
	}
	return nil
}

func (s storer) Export(name, oldpass, transferpass string) ([]byte, error) {
	key, _, err := s.es.Get(name, oldpass)
	if err != nil {
		return nil, err
	}

	res, err := s.es.coder.Encrypt(key, transferpass)
	return res, err
}

func (s storer) Import(name, newpass, transferpass string, data []byte) error {
	key, err := s.es.coder.Decrypt(data, transferpass)
	if err != nil {
		return err
	}

	return s.es.Put(name, newpass, key)
}

func (s storer) Delete(name string) error {
	return s.es.Delete(name)
}

func (s storer) Update(name, oldpass, newpass string) error {
	key, _, err := s.es.Get(name, oldpass)
	if err != nil {
		return err
	}

	// we must delete first, as Putting over an existing name returns an error
	s.Delete(name)

	return s.es.Put(name, newpass, key)
}
