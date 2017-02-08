package lightclient

import crypto "github.com/tendermint/go-crypto"

// KeyInfo is the public information about a key
type KeyInfo struct {
	Name   string
	PubKey crypto.PubKey
}

// Signable represents any transaction we wish to send to tendermint core
// These methods allow us to sign arbitrary Tx with the KeyStore
type Signable interface {
	// Bytes is the immutable data, which needs to be signed
	Bytes() []byte

	// Sign will add a signature and pubkey.
	//
	// Depending on the Signable, one may be able to call this multiple times for multisig
	// Returns error if called with invalid data or too many times
	Sign(pubkey crypto.PubKey, sig crypto.Signature) error

	// SignedBy will return the public key(s) that signed if the signature
	// is valid, or an error if there is any issue with the signature,
	// including if there are no signatures
	SignedBy() ([]crypto.PubKey, error)

	// Signed returns the transaction data as well as all signatures
	// It should return an error if Sign was never called
	SignedBytes() ([]byte, error)
}

// Signer allows one to use a keystore
type Signer interface {
	// Get(name string) (KeyInfo, error)
	Sign(name, passphrase string, tx Signable) error
}

// Poster combines KeyStore and Node to process a Signable and deliver it to tendermint
// returning the results from the tendermint node, once the transaction is processed
// only handles single signatures
type Poster struct {
	server Broadcaster
	signer Signer
}

func NewPoster(server Broadcaster, signer Signer) Poster {
	return Poster{server, signer}
}

// Post will sign the transaction with the given credentials and push it to
// the tendermint server
func (p Poster) Post(sign Signable, keyname, passphrase string) (res TmBroadcastResult, err error) {
	var signed []byte

	err = p.signer.Sign(keyname, passphrase, sign)
	if err != nil {
		return
	}

	signed, err = sign.SignedBytes()
	if err != nil {
		return
	}

	res, err = p.server.Broadcast(signed)
	return
}
