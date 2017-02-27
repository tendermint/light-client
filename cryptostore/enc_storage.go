package cryptostore

import (
	crypto "github.com/tendermint/go-crypto"
	lightclient "github.com/tendermint/light-client"
)

// encryptedStorage needs passphrase to get private keys
type encryptedStorage struct {
	coder Encoder
	store lightclient.Storage
}

func (es encryptedStorage) Put(name, pass string, key crypto.PrivKey) error {
	secret, err := es.coder.Encrypt(key, pass)
	if err != nil {
		return err
	}

	ki := info(name, key)
	return es.store.Put(name, secret, ki)
}

func (es encryptedStorage) Get(name, pass string) (crypto.PrivKey, lightclient.KeyInfo, error) {
	secret, info, err := es.store.Get(name)
	if err != nil {
		return nil, info, err
	}
	key, err := es.coder.Decrypt(secret, pass)
	return key, info, err
}

func (es encryptedStorage) List() ([]lightclient.KeyInfo, error) {
	return es.store.List()
}

func (es encryptedStorage) Delete(name string) error {
	return es.store.Delete(name)
}

// info hardcodes the encoding of keys
func info(name string, key crypto.PrivKey) lightclient.KeyInfo {
	return lightclient.KeyInfo{
		Name:   name,
		PubKey: key.PubKey(),
	}
}
