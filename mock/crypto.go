package mock

import lightclient "github.com/tendermint/light-client"

// PubKey lets us wrap some bytes to provide crypto PubKey for those
// methods that just act on it.
type PubKey struct {
	Val []byte
}

func (p PubKey) assertPubKey() lightclient.PubKey {
	return p
}

func (p PubKey) Bytes() []byte {
	return p.Val
}

// Address is just the pubkey with some constant prepended
func (p PubKey) Address() []byte {
	return append([]byte("Addr:"), p.Val...)
}

func LoadPubKey(data []byte) (PubKey, error) {
	return PubKey{Val: data}, nil
}

// Signature lets us wrap some bytes to provide crypto signature for those
// methods that just act on it.
type Signature struct {
	Val []byte
}

func (s Signature) assertPubKey() lightclient.Signature {
	return s
}

func (s Signature) Bytes() []byte {
	return s.Val
}

func (s Signature) IsZero() bool {
	return len(s.Val) == 0
}
