package mock

import (
	"bytes"
	"encoding/hex"

	crypto "github.com/tendermint/go-crypto"
)

// PubKey lets us wrap some bytes to provide crypto PubKey for those
// methods that just act on it.
type PubKey struct {
	Val []byte
}

func (p PubKey) assertPubKey() crypto.PubKey {
	return p
}

func (p PubKey) Bytes() []byte {
	return p.Val
}

// Address is just the pubkey with some constant prepended
func (p PubKey) Address() []byte {
	return append([]byte("Addr:"), p.Val...)
}

func (p PubKey) KeyString() string {
	return hex.EncodeToString(p.Val)
}

func (p PubKey) VerifyBytes(msg []byte, sig crypto.Signature) bool {
	// TODO
	return true
}

func (p PubKey) Equals(pk crypto.PubKey) bool {
	if _, ok := pk.(PubKey); !ok {
		return false
	}
	return bytes.Equal(p.Bytes(), pk.Bytes())
}

func LoadPubKey(data []byte) (PubKey, error) {
	return PubKey{Val: data}, nil
}

// Signature lets us wrap some bytes to provide crypto signature for those
// methods that just act on it.
type Signature struct {
	Val []byte
}

func (s Signature) assertPubKey() crypto.Signature {
	return s
}

func (s Signature) Bytes() []byte {
	return s.Val
}

func (s Signature) IsZero() bool {
	return len(s.Val) == 0
}

func (s Signature) String() string {
	return hex.EncodeToString(s.Val)
}

func (s Signature) Equals(cs crypto.Signature) bool {
	if _, ok := cs.(Signature); !ok {
		return false
	}
	return bytes.Equal(s.Bytes(), cs.Bytes())
}
