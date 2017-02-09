package tx

import (
	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	lightclient "github.com/tendermint/light-client"
)

// OneSig lets us wrap arbitrary data with a go-crypto signature
//
// TODO: rethink how we want to integrate this with KeyStore so it makes
// more sense (particularly the verify method)
type OneSig struct {
	data   []byte
	sig    crypto.Signature
	pubkey crypto.PubKey
}

func New(data []byte) *OneSig {
	return &OneSig{data: data}
}

func Load(serialized []byte) (*OneSig, error) {
	var s OneSig
	err := wire.ReadBinaryBytes(serialized, &s)
	return &s, err
}

// assertSignable is just to make sure we stay in sync with the Signable interface
func (s *OneSig) assertSignable() lightclient.Signable {
	return s
}

// Bytes returns the original data passed into `NewSig`
func (s *OneSig) Bytes() []byte {
	return s.data
}

// Sign will add a signature and pubkey.
//
// Depending on the Signable, one may be able to call this multiple times for multisig
// Returns error if called with invalid data or too many times
func (s *OneSig) Sign(pubkey crypto.PubKey, sig crypto.Signature) error {
	if pubkey == nil || sig == nil {
		return errors.New("Signature or Key missing")
	}
	if s.sig != nil {
		return errors.New("Transaction can only be signed once")
	}

	// set the value once we are happy
	s.pubkey = pubkey
	s.sig = sig

	return nil
}

// SignedBy will return the public key(s) that signed if the signature
// is valid, or an error if there is any issue with the signature,
// including if there are no signatures
func (s *OneSig) SignedBy() ([]crypto.PubKey, error) {
	if s.pubkey == nil || s.sig == nil {
		return nil, errors.New("Never signed")
	}

	if !s.pubkey.VerifyBytes(s.data, s.sig) {
		return nil, errors.New("Signature doesn't match")
	}

	return []crypto.PubKey{s.pubkey}, nil
}

// SignedBytes serializes the Sig to send it to a tendermint app.
// It returns an error if the Sig was never Signed.
func (s *OneSig) SignedBytes() ([]byte, error) {
	if s.sig == nil {
		return nil, errors.New("Transaction was never signed")
	}
	return wire.BinaryBytes(wrapper{s}), nil
}
