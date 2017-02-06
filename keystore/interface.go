package keystore

import (
	crypto "github.com/tendermint/go-crypto"
	lightclient "github.com/tendermint/light-client"
)

// Encoder is used to encrypt any key with a passphrase for storage
type Encoder interface {
	Encrypt(key crypto.PrivKey, pass string) ([]byte, error)
	Decrypt(data []byte, pass string) (crypto.PrivKey, error)
}

// Storage has many implementation, based on security and sharing requirements
// like disk-backed, mem-backed, vault, db, etc.
type Storage interface {
	Put(name string, key []byte, info lightclient.KeyInfo) error
	Get(name string) ([]byte, lightclient.KeyInfo, error)
	List() ([]lightclient.KeyInfo, error)
	Delete(name string) error
}

// Info hardcodes the encoding of keys
func Info(name string, key crypto.PrivKey) lightclient.KeyInfo {
	pub := key.PubKey()
	return lightclient.KeyInfo{
		Name:    name,
		PubKey:  pub.Bytes(),
		Address: pub.Address(),
	}
}
