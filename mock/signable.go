package mock

import (
	"errors"

	lightclient "github.com/tendermint/light-client"
)

// OneSign is a Signable implementation that can be used to
// record the values and inspect them later.  It performs no validation.
type OneSign struct {
	Data   []byte
	PubKey lightclient.PubKey
	Sig    lightclient.Signature
}

func (o *OneSign) assertSignable() lightclient.Signable {
	return o
}

func (o *OneSign) Bytes() []byte {
	return o.Data
}

func (o *OneSign) Sign(pubkey lightclient.PubKey, sig lightclient.Signature) error {
	if o.PubKey != nil {
		return errors.New("OneSign already signed")
	}
	o.PubKey = pubkey
	o.Sig = sig
	return nil
}

func (o *OneSign) SignedBy() ([]lightclient.PubKey, error) {
	if o.PubKey == nil {
		return nil, errors.New("OneSign never signed")
	}
	return []lightclient.PubKey{o.PubKey}, nil
}

func (o *OneSign) SignedBytes() ([]byte, error) {
	return nil, errors.New("SignedBytes not implemented")
}
