package mock

import (
	"errors"

	lightclient "github.com/tendermint/light-client"
)

// OneSig is a Signable implementation that can be used to
// record the values and inspect them later.  It performs no validation.
type OneSig struct {
	Data   []byte
	PubKey lightclient.PubKey
	Sig    lightclient.Signature
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

func (o *OneSig) Sign(pubkey lightclient.PubKey, sig lightclient.Signature) error {
	if o.PubKey != nil {
		return errors.New("OneSig already signed")
	}
	o.PubKey = pubkey
	o.Sig = sig
	return nil
}

func (o *OneSig) SignedBy() ([]lightclient.PubKey, error) {
	if o.PubKey == nil {
		return nil, errors.New("OneSig never signed")
	}
	return []lightclient.PubKey{o.PubKey}, nil
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
	pubkey lightclient.PubKey
	sig    lightclient.Signature
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

func (m *MultiSig) Sign(pubkey lightclient.PubKey, sig lightclient.Signature) error {
	s := signed{pubkey, sig}
	m.sigs = append(m.sigs, s)
	return nil
}

func (m *MultiSig) SignedBy() ([]lightclient.PubKey, error) {
	if len(m.sigs) == 0 {
		return nil, errors.New("MultiSig never signed")
	}
	keys := make([]lightclient.PubKey, len(m.sigs))
	for i := range m.sigs {
		keys[i] = m.sigs[i].pubkey
	}
	return keys, nil
}

func (m *MultiSig) SignedBytes() ([]byte, error) {
	return nil, errors.New("SignedBytes not implemented")
}
