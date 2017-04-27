package lightclient

import keys "github.com/tendermint/go-crypto/keys"

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
