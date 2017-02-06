package keystore

import (
	crypto "github.com/tendermint/go-crypto"
	lightclient "github.com/tendermint/light-client"
)

// encryptedStorage needs passphrase to get private keys
type encryptedStorage struct {
	coder Encoder
	store Storage
}

func (es encryptedStorage) Put(name, pass string, key crypto.PrivKey) error {
	secret, err := es.coder.Encrypt(key, pass)
	if err != nil {
		return err
	}

	info := Info(name, key)
	return es.store.Put(name, secret, info)
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
