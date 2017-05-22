package lightclient

import (
	crypto "github.com/tendermint/go-crypto"
	keys "github.com/tendermint/go-crypto/keys"
)

type Value interface {
	Bytes() []byte
}

// ValueReader is an abstraction to let us parse application-specific values
type ValueReader interface {
	// ReadValue accepts a key, value pair to decode.  The value bytes must be
	// retained in the returned Value implementation.
	//
	// key *may* be present and can be used as a hint of how to parse the data
	// when your application handles multiple formats
	ReadValue(key, value []byte) (Value, error)
}

type SignableReader interface {
	ReadSignable(data []byte) (keys.Signable, error)
}

// TxReader is anything that can parse incoming transactions and place
// them into a format to send to the node.
//
// Generally the results should implement keys.Signable, in which case they
// will be signed before being posted.  However, there are also cases
// the support unsigned tx (at least dummy, counter, merkleeyes...), which
// can return an implementation of Value.
//
// If this returns anything that doesn't implement either keys.Signable, nor
// Value, then it is considered an error
//
// The signer's pubkey (if any) is passed in, so we can enhace the tx
type TxReader interface {
	// this reads a given json input
	ReadTxJSON([]byte, crypto.PubKey) (interface{}, error)
	// this uses
	ReadTxFlags(interface{}, crypto.PubKey) (interface{}, error)
}
