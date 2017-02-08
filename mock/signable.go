package mock

import (
	"errors"

	crypto "github.com/tendermint/go-crypto"
	lightclient "github.com/tendermint/light-client"
)

// OneSig is a Signable implementation that can be used to
// record the values and inspect them later.  It performs no validation.
type OneSig struct {
	Data   []byte
	PubKey crypto.PubKey
	Sig    crypto.Signature
}

func NewSig(data []byte) *OneSig {
	return &OneSig{Data: data}
}

func (o *OneSig) assertSignable() lightclient.Signable {
	return o
}

func (o *OneSig) Bytes() []byte {
	return o.Data
}

func (o *OneSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error {
	if o.PubKey != nil {
		return errors.New("OneSig already signed")
	}
	o.PubKey = pubkey
	o.Sig = sig
	return nil
}

func (o *OneSig) SignedBy() ([]crypto.PubKey, error) {
	if o.PubKey == nil {
		return nil, errors.New("OneSig never signed")
	}
	return []crypto.PubKey{o.PubKey}, nil
}

func (o *OneSig) SignedBytes() ([]byte, error) {
	return nil, errors.New("SignedBytes not implemented")
}

// MultiSig is a Signable implementation that can be used to
// record the values and inspect them later.  It performs no validation.
//
// It supports an arbitrary number of signatures
type MultiSig struct {
	Data []byte
	sigs []signed
}

type signed struct {
	pubkey crypto.PubKey
	sig    crypto.Signature
}

func NewMultiSig(data []byte) *MultiSig {
	return &MultiSig{Data: data}
}

func (m *MultiSig) assertSignable() lightclient.Signable {
	return m
}

func (m *MultiSig) Bytes() []byte {
	return m.Data
}

func (m *MultiSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error {
	s := signed{pubkey, sig}
	m.sigs = append(m.sigs, s)
	return nil
}

func (m *MultiSig) SignedBy() ([]crypto.PubKey, error) {
	if len(m.sigs) == 0 {
		return nil, errors.New("MultiSig never signed")
	}
	keys := make([]crypto.PubKey, len(m.sigs))
	for i := range m.sigs {
		keys[i] = m.sigs[i].pubkey
	}
	return keys, nil
}

func (m *MultiSig) SignedBytes() ([]byte, error) {
	return nil, errors.New("SignedBytes not implemented")
}
