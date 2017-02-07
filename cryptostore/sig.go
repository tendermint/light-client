package cryptostore

import (
	"github.com/pkg/errors"
	crypto "github.com/tendermint/go-crypto"
	wire "github.com/tendermint/go-wire"
	lightclient "github.com/tendermint/light-client"
)

func NewSig(data []byte) *Sig {
	return &Sig{data: data}
}

// Sig lets us wrap arbitrary data with a go-crypto signature
//
// TODO: rethink how we want to integrate this with KeyStore so it makes
// more sense (particularly the verify method)
type Sig struct {
	data   []byte
	sig    crypto.Signature
	pubkey crypto.PubKey
}

// assertSignable is just to make sure we stay in sync with the Signable interface
func (s *Sig) assertSignable() lightclient.Signable {
	return s
}

// Bytes returns the original data passed into `NewSig`
func (s *Sig) Bytes() []byte {
	return s.data
}

// Sign will add a signature and pubkey.
//
// Depending on the Signable, one may be able to call this multiple times for multisig
// Returns error if called with invalid data or too many times
func (s *Sig) Sign(pubkey lightclient.PubKey, sig lightclient.Signature) error {
	if pubkey == nil || sig == nil {
		return errors.New("Signature or Key missing")
	}
	if s.sig != nil {
		return errors.New("Transaction can only be signed once")
	}

	// make sure the types are truly compatible
	cpk, ok := pubkey.(crypto.PubKey)
	if !ok {
		return errors.New("pubkey must be crypto.PubKey")
	}
	csig, ok := sig.(crypto.Signature)
	if !ok {
		return errors.New("sig must be crypto.Signature")
	}

	// set the value once we are happy
	s.pubkey = cpk
	s.sig = csig

	return nil
}

// SignedBy will return the public key(s) that signed if the signature
// is valid, or an error if there is any issue with the signature,
// including if there are no signatures
func (s *Sig) SignedBy() ([]lightclient.PubKey, error) {
	if s.pubkey == nil || s.sig == nil {
		return nil, errors.New("Never signed")
	}

	if !s.pubkey.VerifyBytes(s.data, s.sig) {
		return nil, errors.New("Signature doesn't match")
	}

	return []lightclient.PubKey{s.pubkey}, nil
}

// SignedBytes serializes the Sig to send it to a tendermint app.
// It returns an error if the Sig was never Signed.
func (s *Sig) SignedBytes() ([]byte, error) {
	if s.sig == nil {
		return nil, errors.New("Transaction was never signed")
	}
	return wire.BinaryBytes(s), nil
}
