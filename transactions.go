package lightclient

// KeyInfo is the public information about a key
type KeyInfo struct {
	Name    string
	Address []byte
	PubKey  []byte
}

// KeyStore represents secure storage of tendermint private keys
// The implementation can specify what types of keys, generally ed25519....
// The implementation is also responsible for deciding how to persist to disk
type KeyStore interface {
	Create(name, passphrase string) error
	List() ([]KeyInfo, error)
	Get(name string) (KeyInfo, error)
	Signature(name, passphrase string, data []byte) ([]byte, error)
	// TODO: import, export, and update (changing passphrase)
}

// Tx represents any transaction we wish to send to tendermint core
// These methods allow us to sign arbitrary Tx with the KeyStore
// TODO: Look at tendermint/types/signable.go
type Tx interface {
	// Bytes is the immutable data, which needs to be signed
	Bytes() []byte

	// AddSignature will add a signature (and address or pubkey as desired)
	// Depending on the Tx, one may be able to call this multiple times for multisig
	// Returns error if called with invalid data or too many times
	Sign(addr, pubkey, sig []byte) error

	// Signed returns bytes ready to post to tendermint
	// It should return an error if AddSignature was never called
	Signed() ([]byte, error)
}

// Poster combines KeyStore and Node to process a Tx and deliver it to tendermint
// returning the results from the tendermint node, once the transaction is processed
// only handles single signatures
type Poster interface {
	Post(tx Tx, keyname, passphrase string) (BroadcastResult, error)
}

// TODO: move this into a subpackage????
type poster struct {
	node Node
	keys KeyStore
}

func NewPoster(node Node, keys KeyStore) Poster {
	return poster{node, keys}
}

func (p poster) Post(tx Tx, keyname, passphrase string) (res BroadcastResult, err error) {
	var info KeyInfo
	var data, sig, signed []byte

	info, err = p.keys.Get(keyname)
	if err != nil {
		return
	}

	data = tx.Bytes()
	sig, err = p.keys.Signature(keyname, passphrase, data)
	if err != nil {
		return
	}

	err = tx.Sign(info.Address, info.PubKey, sig)
	if err != nil {
		return
	}

	signed, err = tx.Signed()
	if err != nil {
		return
	}

	res, err = p.node.Broadcast(signed)
	return
}
